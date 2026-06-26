-- 1.Init_Partition.sql
USE [OK2SHIP_SMT];
GO

-- Kiểm tra và tạo Function trước
IF NOT EXISTS (SELECT 1 FROM sys.partition_functions WHERE name = 'PF_Quarterly')
BEGIN
    CREATE PARTITION FUNCTION PF_Quarterly (DATETIME)
    AS RANGE RIGHT FOR VALUES (
        '2024-01-01', '2024-04-01', '2024-07-01', '2024-10-01',
        '2025-01-01', '2025-04-01', '2025-07-01', '2025-10-01',
        '2026-01-01', '2026-04-01', '2026-07-01', '2026-10-01'
    );
END
GO 

-- Sau đó mới tạo Scheme
IF NOT EXISTS (SELECT 1 FROM sys.partition_schemes WHERE name = 'PS_Quarterly')
BEGIN
    CREATE PARTITION SCHEME PS_Quarterly
    AS PARTITION PF_Quarterly
    ALL TO ([PRIMARY]);
END
GO