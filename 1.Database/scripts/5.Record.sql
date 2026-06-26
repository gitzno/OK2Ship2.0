-- 1. Create table
CREATE TABLE [dbo].[Record] (
    [ID] INT IDENTITY(1,1) NOT NULL,
    [ItemCode] VARCHAR(50) NOT NULL,
    [LotNo] VARCHAR(50) NOT NULL,
    [Type] VARCHAR(50) NOT NULL,
    [Category] VARCHAR(50) NOT NULL,
    [DataLogfile] NVARCHAR(MAX) NULL,
    [DataUser] NVARCHAR(MAX) NULL,
    [CreatedDate] DATETIME DEFAULT GETDATE() NOT NULL,
    [LastModified] DATETIME DEFAULT GETDATE() NOT NULL, 
    CONSTRAINT [PK_Record] PRIMARY KEY CLUSTERED ([ID] ASC)
);
GO

-- 2. Tạo Unique Index 
CREATE UNIQUE NONCLUSTERED INDEX [UQ_Record_4Fields] 
ON [dbo].[Record] (
    [ItemCode] ASC,
    [LotNo] ASC,
    [Type] ASC,
    [Category] ASC
);
GO

-- 3. Trigger update LastModified 
CREATE TRIGGER [dbo].[trg_Record_UpdateLastModified]
ON [dbo].[Record]
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE t
    SET t.[LastModified] = GETDATE()
    FROM [dbo].[Record] t
    INNER JOIN inserted i ON t.[ID] = i.[ID];
END;
GO