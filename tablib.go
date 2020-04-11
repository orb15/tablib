package tablib

import (
	"sync"
	"tablib/tableresult"
	"tablib/validate"
)

//TableRepository is a TableRepository for now
type TableRepository interface {
	AddLuaScript(scriptName, luaScript string) error
	AddTable(yamlBytes []byte) (*validate.ValidationResult, error)
	Execute(scriptName string, callback ParamSpecificationRequestCallback) map[string]string
	List(objectName string) (string, error)
	Pick(tableName string, count int) *tableresult.TableResult
	Roll(tableName string, count int) *tableresult.TableResult
	Search(namePredicate string, tags []string) []*SearchResult
}

//SearchResult holds information about each object discovered during a search
type SearchResult struct {
	Name string
	Type string
	Tags []string
}

//NewTableRepository does what it says on the tin
func NewTableRepository() TableRepository {
	return &concreteTableRepo{
		tableStore:  make(map[string]*tableData),
		scriptStore: make(map[string]*scriptData),
		lock:        &sync.RWMutex{},
	}
}
