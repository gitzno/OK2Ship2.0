package core

/*
	core Table

dataTable.Columns.Add("ID", typeof(int));
dataTable.Columns.Add("ItemCode", typeof(string));
dataTable.Columns.Add("LotNo", typeof(string));
dataTable.Columns.Add("Sheet", typeof(string));
dataTable.Columns.Add("Region", typeof(string));
dataTable.Columns.Add("Sample", typeof(string));
dataTable.Columns.Add("Image", typeof(Image));
dataTable.Columns.Add("Graph", typeof(Image));
dataTable.Columns.Add("Data", typeof(string));
dataTable.Columns.Add("Mode 1: Solder joint crack", typeof(string));
dataTable.Columns.Add("Mode 2: Pad lift", typeof(string));
dataTable.Columns.Add("Mode 3: Solder joint lift", typeof(string));
dataTable.Columns.Add("Mode 4: Intermetallic break", typeof(string));
dataTable.Columns.Add("Mode 5: Component damage", typeof(string));
dataTable.Columns.Add("Mode 6: Component detached", typeof(string));
dataTable.Columns.Add("Mode 7: Flex torn", typeof(string));
dataTable.Columns.Add("Operator", typeof(string));
dataTable.Columns.Add("Time_Update", typeof(string));
dataTable.Columns.Add("Remark", typeof(string));
*/

// Đọc dữ liệu và ghi nó lại
type PeelTest_record struct {
	ID       int
	ItemCode string
	LotNo    string
	Sheet    string
	Region   string
	Sample   string
	Image    string
	Graph    string
	Data     float32
}
