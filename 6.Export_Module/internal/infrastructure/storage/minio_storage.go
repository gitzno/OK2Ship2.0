package storage

import (
	"bytes"
	"context"
	"errors"
	"export_module_service/internal/domain"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
}

func (m *MinioStorage) CopyFile(ctx context.Context, srcPath string, destPath string) error {
	// 1. Phân tách đường dẫn để lấy Bucket và Object
	// Giả định đường dẫn dạng: "bucket-name/folder/file.ext"
	srcParts := strings.SplitN(srcPath, "/", 2)
	destParts := strings.SplitN(destPath, "/", 2)

	if len(srcParts) < 2 || len(destParts) < 2 {
		return fmt.Errorf("đường dẫn không hợp lệ, yêu cầu: 'bucket/path/to/file'")
	}

	srcBucket, srcObject := srcParts[0], srcParts[1]
	destBucket, destObject := destParts[0], destParts[1]

	// 2. Cấu hình nguồn cho lệnh copy
	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcObject,
	}

	// 3. Cấu hình đích cho lệnh copy
	dst := minio.CopyDestOptions{
		Bucket: destBucket,
		Object: destObject,
	}

	// 4. Thực thi lệnh copy phía server
	_, err := m.client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("lỗi khi copy file từ %s sang %s: %w", srcPath, destPath, err)
	}

	return nil
}

func (m *MinioStorage) UploadRawFile(ctx context.Context, folderPath string, objectName string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	bucketName, prefix, err := handlePathFile(folderPath)
	fullPath := prefix + "/" + objectName
	_, err = m.client.PutObject(ctx, bucketName, fullPath, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return fullPath, err
}

func NewMinioStorage(client *minio.Client) domain.LogStorage {
	return &MinioStorage{client: client}
}

// handlePathFile tách filePath thành bucket và prefix an toàn
func handlePathFile(filePath string) (string, string, error) {
	// Loại bỏ dấu "/" ở đầu chuỗi (nếu có) để tránh mảng bị rỗng phần tử đầu
	filePath = strings.TrimPrefix(filePath, "/")

	// Dùng SplitN cắt chuỗi ĐÚNG 1 LẦN tại dấu "/" đầu tiên.
	// Ví dụ: "my-bucket/folder/file.txt" -> ["my-bucket", "folder/file.txt"]
	parts := strings.SplitN(filePath, "/", 2)

	if len(parts) < 2 {
		return "", "", errors.New("đường dẫn không hợp lệ: phải bao gồm bucket và prefix")
	}

	bucket := parts[0]
	prefix := parts[1]

	return bucket, prefix, nil
}

// ListFileLocation lấy danh sách file từ MinIO
func (m *MinioStorage) ListFileLocation(ctx context.Context, filePath string) ([]string, error) {
	// 1. Tách bucket và prefix
	bucket, prefix, err := handlePathFile(filePath + "/")
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	// 2. Cấu hình tùy chọn
	objectOption := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	// 3. Gọi API MinIO
	listObject := m.client.ListObjects(ctx, bucket, objectOption)
	var result []string

	// 4. Duyệt kết quả và kiểm tra lỗi chặt chẽ
	for object := range listObject {
		// BẮT BUỘC: Kiểm tra lỗi trong quá trình stream dữ liệu
		if object.Err != nil {
			return nil, fmt.Errorf("lỗi khi đọc object từ minio: %v", object.Err)
		}

		result = append(result, bucket+"/"+object.Key)
	}

	return result, nil
}

func (m *MinioStorage) GetLogStream(ctx context.Context, filePath string) (io.ReadCloser, error) {
	bucket, objectName, err := handlePathFile(filePath)
	if err != nil {
		return nil, err
	}

	object, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("lỗi khởi tạo stream lấy file: %v", err)
	}

	_, err = object.Stat()
	if err != nil {
		object.Close()
		return nil, fmt.Errorf("file không tồn tại hoặc không thể truy cập: %v", err)
	}

	// Trả về object (chính là một io.ReadCloser)
	return object, nil
}

// MoveFile thực hiện copy file sang vị trí mới (kèm tên mới) và xóa file cũ
func (m *MinioStorage) MoveFile(ctx context.Context, srcPath string, destPath string) error {
	// 1. Phân tách đường dẫn Nguồn
	srcParts := strings.SplitN(srcPath, "/", 2)
	if len(srcParts) < 2 {
		return fmt.Errorf("đường dẫn nguồn không hợp lệ: %s", srcPath)
	}
	srcBucket, srcObject := srcParts[0], srcParts[1]

	// 2. Phân tách đường dẫn Đích (Chứa tên mới)
	destParts := strings.SplitN(destPath, "/", 2)
	if len(destParts) < 2 {
		return fmt.Errorf("đường dẫn đích không hợp lệ: %s", destPath)
	}
	destBucket, destObject := destParts[0], destParts[1]

	// 3. Cấu hình Copy
	srcOpts := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcObject,
	}
	destOpts := minio.CopyDestOptions{
		Bucket: destBucket,
		Object: destObject,
	}

	// 4. Thực hiện Copy sang tên/vị trí mới
	_, err := m.client.CopyObject(ctx, destOpts, srcOpts)
	if err != nil {
		return fmt.Errorf("lỗi khi copy file từ %s sang %s: %w", srcPath, destPath, err)
	}

	// 5. Xóa file ở vị trí cũ
	err = m.client.RemoveObject(ctx, srcBucket, srcObject, minio.RemoveObjectOptions{})
	if err != nil {
		// Lưu ý: Nếu copy thành công nhưng xóa thất bại, bạn có thể bị rác dữ liệu (file tồn tại ở cả 2 nơi)
		return fmt.Errorf("đã copy thành công nhưng lỗi khi xóa file gốc %s: %w", srcPath, err)
	}

	return nil
}

// MoveBatchFiles nhận vào một map: Key là đường dẫn cũ, Value là đường dẫn mới
func (m *MinioStorage) MoveBatchFiles(ctx context.Context, fileMap map[string]string) error {
	for src, dest := range fileMap {
		if err := m.MoveFile(ctx, src, dest); err != nil {
			// Tùy theo logic nghiệp vụ: bạn có thể log lỗi rồi chạy tiếp, hoặc return lỗi để dừng ngay
			return fmt.Errorf("tiến trình di chuyển hàng loạt bị dừng tại file %s: %w", src, err)
		}
	}
	return nil
}
