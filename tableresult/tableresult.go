package tableresult

//TableResult holds the final result of a table run
type TableResult struct {
	Result []string
	Data   map[string]string
	Log    []string
}

//NewTableResult does what it says on the tin
func NewTableResult() *TableResult {
	tr := &TableResult{
		Result: make([]string, 0, 1),
		Data:   make(map[string]string),
		Log:    make([]string, 0, 1),
	}
	return tr
}

//AddLog adds the string to the log
func (tr *TableResult) AddLog(msg string) {
	tr.Log = append(tr.Log, msg)
}

//AddResult adds the string to the results (for multi-rolls)
func (tr *TableResult) AddResult(msg string) {
	tr.Result = append(tr.Result, msg)
}
