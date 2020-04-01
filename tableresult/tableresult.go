package tableresult

//TableResult holds the final result of a table run
type TableResult struct {
	Result string
	Log    []string
}

//NewTableResult does what it says on the tin
func NewTableResult() *TableResult {
	tr := &TableResult{
		Result: "",
		Log:    make([]string, 0, 1),
	}
	return tr
}
