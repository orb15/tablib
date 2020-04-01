package tsm

//TableStateMap holds a parse tree for table expansion
type TableStateMap map[string]string

//NewTableStateMap does what it says on the tin
func NewTableStateMap() TableStateMap {
	return make(TableStateMap)
}

//ExpansionInfo holds data about the next string expansion to process
type ExpansionInfo struct {
	ToExpand string
	Key      string
	Done     bool
}

//FindNextForExpansion searches the TSM for a table reference to expand
func (tsm *TableStateMap) FindNextForExpansion() *ExpansionInfo {
	return nil
}
