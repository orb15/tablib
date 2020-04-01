package tablib

import (
	"sync"
	"tablib/util"
)

//TableRepository is a TableRepository for now
type TableRepository interface {
	AddLuaScript(scriptName, luaScript string) error
	AddTable(yamlBytes []byte) (*util.ValidationResult, error)
	Execute(scriptName string) (map[string]string, error)
	List(objectName string) (string, error)
	Pick(tableName string, count int) *TableResult
	Roll(tableName string, count int) *TableResult
	Search(namePredicate string, tags []string) []*SearchResult
}

//SearchResult holds information about each object discovered during a search
type SearchResult struct {
	Name string
	Type string
	Tags []string
}

//TableResult holds the final result of a table run
type TableResult struct {
	Result string
	Log    []string
}

//NewTableRepository does what it says on the tin
func NewTableRepository() TableRepository {
	return &concreteTableRepo{
		tableStore:  make(map[string]*tableData),
		scriptStore: make(map[string]*scriptData),
		lock:        &sync.RWMutex{},
	}
}
