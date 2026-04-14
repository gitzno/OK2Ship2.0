-- 2.Select_Super.sql
USE [OK2SHIP_SMT];
GO

CREATE OR ALTER PROCEDURE [dbo].[sp_UniversalSearch]
    @TableName NVARCHAR(128),
    @ItemCode VARCHAR(50),
    @LotNo VARCHAR(50),
    @Type VARCHAR(20)
AS
BEGIN
    SET NOCOUNT ON;

    IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = @TableName)
    BEGIN
        RAISERROR('Tên bảng không tồn tại!', 16, 1);
        RETURN;
    END

    DECLARE @SafeTableName NVARCHAR(256) = QUOTENAME(@TableName);
    DECLARE @DynamicSQL NVARCHAR(MAX);

    -- Lắp ráp SQL: Tìm ID trên Index (Seek) -> Join ngược lại lấy Data (Lookup)
    SET @DynamicSQL = N'
        WITH FoundID AS (
            -- Bước 1: Tìm nhanh trong 3 tháng gần nhất
            SELECT TOP 1 ID, CreatedDate
            FROM dbo.' + @SafeTableName + N'
            WHERE ItemCode = @p_ItemCode AND LotNo = @p_LotNo AND Type = @p_Type
              AND CreatedDate >= DATEADD(MONTH, -3, GETDATE())
            
            UNION ALL
            
            -- Bước 2: Nếu bước 1 không có (0 dòng), mới tìm toàn bộ lịch sử
            SELECT TOP 1 ID, CreatedDate
            FROM dbo.' + @SafeTableName + N'
            WHERE ItemCode = @p_ItemCode AND LotNo = @p_LotNo AND Type = @p_Type
            AND NOT EXISTS (
                 SELECT 1 FROM dbo.' + @SafeTableName + N' 
                 WHERE ItemCode = @p_ItemCode AND LotNo = @p_LotNo AND Type = @p_Type
                 AND CreatedDate >= DATEADD(MONTH, -3, GETDATE())
            )
        )
        -- Bước 3: Join lấy full dữ liệu (Data, LocationImg)
        SELECT t.* FROM dbo.' + @SafeTableName + N' t
        INNER JOIN FoundID f ON t.ID = f.ID AND t.CreatedDate = f.CreatedDate;
    ';

    EXEC sp_executesql 
        @stmt = @DynamicSQL,
        @params = N'@p_ItemCode VARCHAR(50), @p_LotNo VARCHAR(50), @p_Type VARCHAR(20)',
        @p_ItemCode = @ItemCode,
        @p_LotNo = @LotNo,
        @p_Type = @Type;
END;
GO