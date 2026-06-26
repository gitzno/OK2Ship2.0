package usecase

import (
	"context"
	"encoding/json"
	"export_module_service/internal/domain"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

type PeelTestUC struct {
	storage    domain.LogStorage ``
	repository domain.RecordRepository
}

func NewPeelTestUC(s domain.LogStorage, repo domain.RecordRepository) *PeelTestUC {
	return &PeelTestUC{storage: s, repository: repo}
}

// padLeft là hàm tiện ích giúp chèn ký tự vào bên trái chuỗi cho đủ độ dài mong muốn.
func padLeft(s string, length int, padStr string) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(padStr, length-len(s)) + s
}

// StandardizeData chuẩn hóa ItemCode (6 ký tự) và LotNo (5 ký tự phần chính, 5 ký tự mở rộng nếu có)
func StandardizeData(itemCode, lotNo string) (string, string) {
	// Chuẩn hóa ItemCode
	stdItemCode := padLeft(itemCode, 6, "0")

	// Chuẩn hóa LotNo
	lotParts := strings.Split(lotNo, "-")
	stdLotBase := padLeft(lotParts[0], 5, "0")

	if len(lotParts) > 1 {
		stdLotExt := padLeft(lotParts[1], 5, "0")
		return stdItemCode, fmt.Sprintf("%s-%s", stdLotBase, stdLotExt)
	}

	return stdItemCode, stdLotBase
}

// VerifyFolderInfo thực hiện trọn vẹn: Chuẩn hóa đầu vào -> Tách chuỗi Folder -> Đối chiếu
func VerifyFolderInfo(folderPath, rawItemCode, rawLotNo string) (bool, string, string) {
	// 1. Chuẩn hóa dữ liệu đầu vào
	expectedItemCode, expectedLotNo := StandardizeData(rawItemCode, rawLotNo)

	// 2. Lấy tên folder lõi
	folderName := folderPath
	if strings.Contains(folderPath, "/") {
		parts := strings.SplitN(folderPath, "/", 2)
		folderName = parts[1]
	}

	// 3. Tính số đoạn cần cắt
	lotNoParts := strings.Split(expectedLotNo, "-")
	numLotNoParts := len(lotNoParts)
	expectedSegmentsCount := 1 + 1 + numLotNoParts + 1

	// 4. Cắt chuỗi
	segments := strings.SplitN(folderName, "-", expectedSegmentsCount)

	// Nếu chuỗi sai định dạng hoàn toàn, trả về false cùng giá trị chuẩn hóa
	if len(segments) < expectedSegmentsCount {
		return false, expectedItemCode, expectedLotNo
	}

	// 5. Trích xuất
	extractedItemCode := segments[1]
	extractedLotNo := strings.Join(segments[2:2+numLotNoParts], "-")

	// 6. Đối chiếu
	isValid := (extractedItemCode == expectedItemCode) && (extractedLotNo == expectedLotNo)

	return isValid, extractedItemCode, extractedLotNo
}

// {
// "file_path": "ok2ship-logfile/6CRMEV-7D0760-00010- WITHOUT SUS-SOLDER SINGAPORE-Peel",
// "item_code": "7D0760",
// "lot_no": "00010",
// "type": "peel_test"
// }
const CATEGORY = "PEEL_TEST"

const IMAGE_BUCKET_NAME = "ok2ship-images"

func extractNumber(filename string) int {
	nameOnly := strings.TrimSuffix(filename, filepath.Ext(filename))
	num, _ := strconv.Atoi(nameOnly)
	return num
}

func GenerateImageRenameMap(filePaths []string) map[string]string {
	// 1. Lọc chỉ lấy file .jpg
	var jpgFiles []string
	for _, path := range filePaths {
		if strings.ToLower(filepath.Ext(path)) == ".jpg" {
			jpgFiles = append(jpgFiles, path)
		}
	}

	// 2. Sắp xếp file theo số thứ tự (dựa trên tên file hiện tại)
	// Giả sử tên file là "750.jpg", "751.jpg"... ta extract số ra để sort
	sort.Slice(jpgFiles, func(i, j int) bool {
		numI := extractNumber(filepath.Base(jpgFiles[i]))
		numJ := extractNumber(filepath.Base(jpgFiles[j]))
		return numI < numJ
	})

	// 3. Tạo map đổi tên
	renameMap := make(map[string]string)
	for i, path := range jpgFiles {
		// i là index (0, 1, 2...), ta cần bắt đầu từ 1 -> i+1
		newName := strconv.Itoa(i+1) + ".jpg"
		renameMap[path] = newName
	}

	return renameMap
}

// 1. Thực thi hàm ImportLogFile (Đọc logfile được người dùng up lên object storage -> lưu DB và ảnh di chuyển đi object storage)
func (uc *PeelTestUC) ImportLogFile(ctx context.Context, fileAddress string, itemCode string, lotNo string) error {

	// kiểm tra itemCode lotno của folder vừa upload
	log.Printf("Bắt đầu xử lý file: %s cho Item: %s, Lot: %s", fileAddress, itemCode, lotNo)
	prime, _, _ := VerifyFolderInfo(fileAddress, itemCode, lotNo)
	if prime == false {
		return fmt.Errorf("itemCode or lotNo is invalid")
	}

	// chuẩn bị đường dẫn lưu trữ ảnh
	folderImagePath := fmt.Sprintf("%s/%s_%s", IMAGE_BUCKET_NAME, CATEGORY, itemCode, lotNo)

	// lấy danh sách file
	listFiles, err := uc.storage.ListFileLocation(ctx, fileAddress)

	// chuẩn bị số thứ tự cho ảnh chụp
	dic := GenerateImageRenameMap(listFiles)

	if err != nil {
		return fmt.Errorf("Lỗi khi lấy danh sách %m", err)
	}

	resultJson, err := BatchProcessFile(ctx, uc.storage, listFiles, folderImagePath, dic)
	fmt.Println(resultJson)
	record := domain.Record{ItemCode: itemCode, LotNo: lotNo, Category: CATEGORY, Type: "NPI", DataLogfile: resultJson}
	err = uc.repository.Upsert(
		ctx,
		record,
	)
	return err
}

func BatchProcessFile(ctx context.Context, storage domain.LogStorage, filePath []string, folderImagePath string, dic map[string]string) (string, error) {
	jobs := make(chan string, len(filePath))
	results := make(chan FileProcessResult, len(filePath))
	var wg sync.WaitGroup

	// khởi tạo 10 worker chạy song song
	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				results <- processSingleFile(ctx, storage, path, folderImagePath, dic)
			}
		}()
	}

	//producer: đẩy việc cho worker
	for _, path := range filePath {
		jobs <- path
	}
	close(jobs)

	//chờ các worker làm xong
	go func() {
		wg.Wait()
		close(results)

	}()
	dic_result := make(map[string]*FileProcessResult)

	// consume nhận kết quả
	for res := range results {
		if res.Err != nil {
			log.Printf("File: %s - Err %s\n", res.FileName, res.Err)
		} else {
			if existing, ok := dic_result[res.key]; !ok {
				// Nếu chưa có, gán bản sao của res vào map
				dic_result[res.key] = &res
			} else {
				// Đã có: Cập nhật các trường nếu trường mới khác rỗng
				if existing.MaxValue == "" {
					dic_result[res.key].MaxValue = res.MaxValue
				}
				if existing.LocationImage == "" {
					dic_result[res.key].LocationImage = res.LocationImage
				}
				if existing.LocationGraph == "" {
					dic_result[res.key].LocationGraph = res.LocationGraph
				}

			}
		}
	}

	// Chuyển đổi map thành JSON
	jsonData, err := json.Marshal(dic_result)
	if err != nil {
		return "", fmt.Errorf("lỗi khi tạo json: %w", err)
	}

	return string(jsonData), nil
}

type FileProcessResult struct {
	key           string `json:"key"`
	FileName      string `json:"file_name"`
	MaxValue      string `json:"max_value"`
	LocationGraph string `json:"location_graph"`
	LocationImage string `json:"location_image"`
	Err           error  `json:"err"`
}

func processSingleFile(ctx context.Context, storage domain.LogStorage, filePath string, imagePath string, dic map[string]string) FileProcessResult {
	result := FileProcessResult{FileName: filePath}

	// lấy tên của file làm số thứ tự của file
	fileName := strings.Split(strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1], ".")
	result.key = fileName[0]
	if _, err := strconv.Atoi(fileName[0]); err != nil {
		result.Err = err
		return result
	}

	// get Stream
	stream, err := storage.GetLogStream(ctx, filePath)
	if err != nil {
		result.Err = err
		return result
	}
	defer stream.Close()
	if fileName[1] == "xlsx" {
		imagePath += "/graph"
		f, err := excelize.OpenReader(stream)
		if err != nil {
			result.Err = err
			return result
		}

		sheet := f.GetSheetName(0)

		coords, _ := f.SearchSheet(sheet, "Max")
		// tìm giá trị max
		if len(coords) > 0 {
			col, row, _ := excelize.CellNameToCoordinates(coords[0])
			rightCell, err := getRightCell(f, sheet, col, row)
			if err == nil {
				result.MaxValue, _ = f.GetCellValue(sheet, rightCell)
				result.MaxValue = strings.Replace(result.MaxValue, " ", "", -1)
			}
		}

		// 2. Tìm ảnh tên "Picture 1"
		pics, _ := f.GetPictureCells(sheet)
		for _, pic := range pics {
			picture, _ := f.GetPictures(sheet, pic)
			if picture[0].Extension == ".png" || picture[0].Extension == ".jpg" || picture[0].Extension == ".jpeg" {
				if picture[0].Format.Name == "Picture 1" {
					result.LocationGraph = picture[0].Format.Name
					// upload ảnh lên đường dẫn đã khởi tạo

					location, err := storage.UploadRawFile(ctx, imagePath, fileName[0]+picture[0].Extension, picture[0].File, "image/"+strings.TrimPrefix(picture[0].Extension, "."))
					if err != nil {
						result.Err = err

					} else {
						result.LocationGraph = location
					}

				}
				break
			}
		}
		return result
	}
	if fileName[1] == "jpg" || fileName[1] == "jpeg" || fileName[1] == "png" {
		// nếu là file ảnh thông thường copy file ảnh sang bên
		newPath := imagePath + "/image/" + dic[filePath]
		result.key = strings.Replace(dic[filePath], "."+fileName[1], "", -1)
		err := storage.CopyFile(ctx, filePath, newPath)
		if err != nil {
			result.Err = err
		} else {
			result.LocationImage = newPath
		}
		return result
	}
	result.Err = fmt.Errorf("File không nằm trong spec")
	return result
}

// getRightCell tự động tìm ô liền kề bên phải, nhảy qua các cột bị merge nếu có
func getRightCell(f *excelize.File, sheetName string, col, row int) (string, error) {
	// Chuyển đổi tọa độ hiện tại thành tên ô (VD: 1,1 -> "A1")
	startCell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return "", err
	}

	nextCol := col + 1

	// Lấy danh sách các vùng đã merge trong sheet
	mergedCells, err := f.GetMergeCells(sheetName)
	if err != nil {
		return "", err
	}

	// Kiểm tra xem ô hiện tại có nằm ở đầu một vùng merge không
	for _, merge := range mergedCells {
		if merge.GetStartAxis() == startCell {
			endAxis := merge.GetEndAxis()
			// Lấy cột của ô kết thúc vùng merge
			endCol, _, err := excelize.CellNameToCoordinates(endAxis)
			if err == nil {
				nextCol = endCol + 1 // Cột bên phải thực sự sẽ nằm sau vùng gộp
			}
			break
		}
	}

	return excelize.CoordinatesToCellName(nextCol, row)
}

// 2. Thực thi hàm ExportData (Đọc từ DB -> Tạo file -> Đẩy lên MinIO)
func (uc *PeelTestUC) ExportData(ctx context.Context, itemCode string, lotNo string) error {
	log.Printf("Bắt đầu export: cho Item: %s, Lot: %s", itemCode, lotNo)

	// A. Lấy dữ liệu từ Database (ví dụ: lấy dữ liệu theo khoảng thời gian)
	// data, err := uc.repo.FetchDataForExport(ctx)

	// B. Tạo file Excel mới trong bộ nhớ
	// file, err := createExcelFile(data)

	// C. Upload file vừa tạo lên MinIO
	// _, err = uc.storage.UploadImage(ctx, fileBytes, ".xlsx")

	return nil
}
