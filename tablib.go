package tablib

import (
	"tablib/table"
	"tablib/util"
)

//TableRepository is a TableRepository for now
type TableRepository interface {
	Add(yamlBytes []byte) (*util.ValidationResult, error)
	List(searchExpr string) ([]*ListResponse, error)
}

//ListResponse is returned on repo List() operation
type ListResponse struct {
	Name string
	Note string
	Type string
}

//NewTableRepository does what it says on the tin
func NewTableRepository() TableRepository {
	return &concreteTableRepo{
		tableData: make(map[string]*table.Table),
	}
}
