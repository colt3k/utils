package data

//Col column data store for Data as a string and Name as a string i.e. Key/Value
type Col struct {
	Data string
	Name string
}

//Row data stored for Columns
type Row struct {
	Cols []Col
}

//Table data stored for Rows
type Table struct {
	Rows []Row
}
