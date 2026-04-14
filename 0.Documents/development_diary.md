# Week 3:
## Database
Ưu điểm của kiến trúc chia mảnh của tôi
```
Action              | Normal                      | Kiến trúc mới
--------------------------------------------------------------------------
Tìm kiếm DATA mới   | Quét toàn bộ dữ liệu        | Tìm quý gần nhất
--------------------------------------------------------------------------
Dung lượng lưu trữ  | 100% Dung lượng thực tế     | Giảm ~60%
--------------------------------------------------------------------------
Tìm kiếm DATA cũ    | Quét toàn bộ dữ liệu        | Ổn định(Olog(n))
--------------------------------------------------------------------------
Xóa dữ liệu hết hạn | tìm kiếm và xóa             | tức thì (O(1))
--------------------------------------------------------------------------
```

## Object storage & API Upload service
Tôi đã phải chuyển từ MINIO sang SeaweedFS do lo ngại vấn đề pháp lý nếu phần mềm này được triển khai.
Minio đã chuyển giấy phép.

- Triển khai việc lưu trữ ví nóng => ví lạnh => Hết hạn
- Triển khai API Upload với golang


# Week 2:
## Database
Như tuần trước tôi sẽ làm lần lượt hệ thông database của tôi, ở đây tôi khởi tạo tự động bằng DOCKER và code T-SQL:
 - tạo docker file cho pull mssql về và chạy entrypoint để tự động chạy các scrpit trong folder ./script 
 - đẩy docker file vào trong docker-compose thêm biến ${MSSQL_PASSWORD} vào trong env
 => docker compose tự động khởi tạo

0.Init_DB.sql: Tạo database OK2SHIP_SMT 2.0 
```
USE [master]
GO
/****** Object:  Database [OK2SHIP_SMT]    Script Date: 3/20/2026 5:31:44 PM ******/
CREATE DATABASE [OK2SHIP_SMT]

```

1.init_partition: 
 - Tạo Partition Function (PF_Quarterly): Cắm các mốc thời gian chia theo từng quý (2024, 2025, 2026...). Đây là "Luật chia".

- Tạo Partition Scheme (PS_Quarterly): Gắn luật chia này lên ổ cứng. Đây là "Bản đồ lưu trữ". 
```
USE [OK2SHIP_SMT];
GO

-- =========================================================
-- GIAI ĐOẠN 1: TẠO LUẬT CHIA (FUNCTION) VÀ BẢN ĐỒ (SCHEME)
-- =========================================================

-- 2. Tạo Partition Function: Chia theo từng Quý
-- Tôi đã tạo sẵn các mốc từ năm 2024 đến hết 2026 cho bạn.
CREATE PARTITION FUNCTION PF_Quarterly (DATETIME)
AS RANGE RIGHT FOR VALUES (
    '2024-01-01', '2024-04-01', '2024-07-01', '2024-10-01',
    '2025-01-01', '2025-04-01', '2025-07-01', '2025-10-01',
    '2026-01-01', '2026-04-01', '2026-07-01', '2026-10-01'
);
GO

-- 3. Tạo Partition Scheme: Gắn Luật chia trên vào ổ đĩa PRIMARY
CREATE PARTITION SCHEME PS_Quarterly
AS PARTITION PF_Quarterly
ALL TO ([PRIMARY]);
GO
```


2. Tạo ra funciton truy vấn đỉnh cao
```
USE [OK2SHIP_SMT];
GO

CREATE PROCEDURE [dbo].[sp_UniversalSearch]
    @TableName NVARCHAR(128),  -- Tham số nhận tên bảng (ví dụ: 'PEEL_TEST_NAS', 'ACF_BONDING_NAS')
    @ItemCode VARCHAR(20),
    @LotNo VARCHAR(20),
    @Type VARCHAR(20)
AS
BEGIN
    SET NOCOUNT ON;

    -- 1. BẢO MẬT: Kiểm tra xem bảng có tồn tại không (Chống SQL Injection)
    IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = @TableName)
    BEGIN
        RAISERROR('Tên bảng không tồn tại trong hệ thống!', 16, 1);
        RETURN;
    END

    -- Hàm QUOTENAME sẽ tự động bọc tên bảng thành dạng [TenBang] cho chuẩn cú pháp
    DECLARE @SafeTableName NVARCHAR(256) = QUOTENAME(@TableName);
    DECLARE @DynamicSQL NVARCHAR(MAX);

    -- 2. "LẮP RÁP" câu lệnh SQL Thác đổ
    SET @DynamicSQL = N'
        -- Bước 1: Ưu tiên tìm 3 tháng gần nhất (Quét ít Partition)
        SELECT * INTO #TempResult 
        FROM dbo.' + @SafeTableName + N'
        WHERE ItemCode = @p_ItemCode AND LotNo = @p_LotNo AND Type = @p_Type
          AND CreatedDate >= DATEADD(MONTH, -3, GETDATE());

        -- Kiểm tra nếu có dữ liệu ở Bước 1
        IF EXISTS (SELECT 1 FROM #TempResult)
        BEGIN
            -- Trả về luôn cho lẹ
            SELECT * FROM #TempResult;
        END
        ELSE
        BEGIN
            -- Bước 2: Không thấy thì mới tìm full toàn bộ lịch sử (Dự phòng)
            SELECT * FROM dbo.' + @SafeTableName + N'
            WHERE ItemCode = @p_ItemCode AND LotNo = @p_LotNo AND Type = @p_Type;
        END
    ';

    -- 3. THỰC THI câu lệnh SQL Động một cách an toàn
    EXEC sp_executesql 
        @stmt = @DynamicSQL,
        @params = N'@p_ItemCode VARCHAR(20), @p_LotNo VARCHAR(20), @p_Type VARCHAR(20)',
        @p_ItemCode = @ItemCode,
        @p_LotNo = @LotNo,
        @p_Type = @Type;

END;
GO

```


3.x Tạo ra các bảng dữ liệu:
Theo cấu trúc bắt buộc phải có:
[ ] Có cột CreatedDate kiểu DATETIME chưa?

[ ] Khóa chính đã có cả ID và CreatedDate chưa?

[ ] Cuối lệnh CREATE TABLE đã có ON PS_Quarterly(CreatedDate) chưa?

[ ] Các cột tìm kiếm đã được gom vào Covering Index chưa?

[ ] Tên các cột lõi có khớp với Stored Procedure vạn năng không?

Form tạo bảng truy vấn nhanh:

```
-- 1. Tên bảng: Thay thế [PEEL_TEST_NAS] bằng tên bảng mới của bạn
CREATE TABLE [dbo].[PEEL_TEST_NAS](
    [ID] [INT] IDENTITY(1,1) NOT NULL,
    
    -- Các cột định danh (Phải khớp tên để Stored Procedure chạy được)
    [ItemCode] [VARCHAR](50) NOT NULL, 
    [LotNo] [VARCHAR](50) NOT NULL,
    [Type] [VARCHAR](20) NOT NULL,
    
    -- Các cột dữ liệu nặng
    [Data] [NVARCHAR](MAX) NULL,             -- Lưu chuỗi JSON kết quả test
    [LocationImg] [NVARCHAR](500) NULL,      -- Lưu đường dẫn ảnh trên MinIO/Server
    
    -- Cột mốc thời gian (Bắt buộc để phân mảnh)
    [CreatedDate] [DATETIME] NOT NULL 
        CONSTRAINT [DF_TableName_CreatedDate] DEFAULT (GETDATE()),

    -- KHÓA CHÍNH: Bắt buộc chứa CreatedDate để phân mảnh
    CONSTRAINT [PK_PEEL_TEST_NAS] PRIMARY KEY CLUSTERED 
    (
        [ID] ASC, 
        [CreatedDate] ASC
    )
) ON [PS_Quarterly]([CreatedDate]); -- Gắn bảng vào Scheme
GO

-- 2. Tên Index: Thay [IX_PEEL_TEST_Search] cho phù hợp với tên bảng
-- Tạo mục lục siêu tốc chứa sẵn dữ liệu nặng (Covering Index)
CREATE NONCLUSTERED INDEX [IX_PEEL_TEST_Search] 
ON [dbo].[PEEL_TEST_NAS] ([ItemCode], [LotNo], [Type])
ON [PS_Quarterly]([CreatedDate]);
GO
```
# Week 1:

## Background

Tập chung vào logic cốt lõi của dự án xuất dữ liệu

Dự án này để chuyển đổi dự án cũ OK2Ship sang một hệ thống mới trước đây chỉ tập chung vào giải quyết vấn đề trước mắt là:

> UI cho công nhân sử dụng => lưu database => xuất file.

> Gói gọn trong một **project winform**.

Nhưng khi dữ liệu ngày càng nhiều và phức tạp thì phần mềm đã giảm hiệu suất cũng như bảo trì gặp nhiều cản trở khó khăn.

=> Vì vậy tôi quyết định chuyển đổi nó thành hệ thống hướng sự kiện thế hệ mới.

> Trong quá trình sử dụng tôi sẽ vừa làm tài liệu lên kiến trúc và phương pháp chuyển đổi sang cho có thể chuyển đổi nhẹ nhàng sang hệ thống mới.

## Database

Database ban đầu được sử dụng hoàn toàn bằng mssql - lưu ảnh - lưu dữ liệu từng cột VD:

```
//Old table ACF_BONDING
CREATE TABLE [dbo].[ACF_BONDING](
  [ID] [int] NOT NULL,
  [ItemCode] [nvarchar](max) NULL,
  [LotNo] [nvarchar](max) NULL,
  [Pcs_No] [nvarchar](max) NULL,
  [Image_Before] [varbinary](max) NULL,
  [Image_After] [varbinary](max) NULL,
  [Graph] [varbinary](max) NULL,
  [Data] [nvarchar](max) NULL,
  [Operator] [nvarchar](max) NULL,
  [Time_Update] [nvarchar](max) NULL,
  [Remark] [nvarchar](max) NULL
)
```

Hiện có 3 vấn đề lớn ở trong bảng dữ liệu này:

- Dữ liệu có kèm dữ liệu phức tạp
- Dữ liệu trên thực tế chỉ truy vấn đến ItemCode và LotNo
- Chưa đánh index dữ liệu



Hệ quả dẫn tới hệ thống bị treo và không sử dụng được. Dựa vào các vấn đề gặp phải tôi sử dụng kiến trúc cơ sở dữ liệu mới như sau:

- Chuyển đổi từ lưu trữ byte[] Database => object storage.
- Chuyển đổi từ lưu 1 dữ liệu cho 1 itemcode lotno nhiều dòng => 1 dòng
- Đánh Index cho những cột thường truy vấn.
- Sử dụng kỹ thuật patition table để phân mảnh phục vụ việc dữ liệu hết hạn và xóa



## Bước đầu Tập chung vào logic cốt lõi của dự án xuất dữ liệu

Tôi chọn golang là ngôn ngữ đảm nhiệm logic phần này.

### Tại sao là GO / GOLANG

- Tại tôi thích

### Thiết lập cấu trúc thư mục theo chuẩn [Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_vi.md)

- Tạo Folder 1.ExportModule
  mkdir cmd\api, cmd\worker, internal\shipment, internal\export, internal\config, pkg\logger, pkg\database, api, deployments
  go mod init ok2ship2.0_export_module

###
