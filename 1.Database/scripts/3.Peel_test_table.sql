-- 3.Peel_test_table.sql
USE [OK2SHIP_SMT];
GO

IF OBJECT_ID('[dbo].[PEEL_TEST_NAS]', 'U') IS NOT NULL
    DROP TABLE [dbo].[PEEL_TEST_NAS];
GO

CREATE TABLE [dbo].[PEEL_TEST_NAS](
    [ID] [INT] IDENTITY(1,1) NOT NULL,
    [ItemCode] [VARCHAR](50) NOT NULL, 
    [LotNo] [VARCHAR](50) NOT NULL,
    [Type] [VARCHAR](20) NOT NULL,
    [Data] [NVARCHAR](MAX) NULL,             
    [LocationImg] [NVARCHAR](500) NULL,      
    [CreatedDate] [DATETIME] NOT NULL 
        CONSTRAINT [DF_PEEL_TEST_NAS_CreatedDate] DEFAULT (GETDATE()),

    CONSTRAINT [PK_PEEL_TEST_NAS] PRIMARY KEY CLUSTERED ([ID] ASC, [CreatedDate] ASC)
) ON [PS_Quarterly]([CreatedDate])
WITH (DATA_COMPRESSION = PAGE); -- Phải có dấu ; ở đây
GO

CREATE UNIQUE NONCLUSTERED INDEX [UX_PEEL_TEST_NAS_Identity] 
ON [dbo].[PEEL_TEST_NAS] ([ItemCode], [LotNo], [Type], [CreatedDate])
ON [PS_Quarterly]([CreatedDate])
WITH (DATA_COMPRESSION = PAGE); -- Phải có dấu ; ở đây
GO