package tablib

import (
	"tablib/table"
	"tablib/util"

	yaml "gopkg.in/yaml.v2"
)

//TableRepo is a TableRepo for now
type TableRepo interface {
	Add(yamlBytes []byte) (*util.ValidationResult, error)
}

//NewTableRepo does what it says on the tin
func NewTableRepo() TableRepo {
	return &concreteTableRepo{
		tableData: make(map[string]*table.Table),
	}
}

type concreteTableRepo struct {
	tableData map[string]*table.Table
}

func (cr *concreteTableRepo) Add(yamlBytes []byte) (*util.ValidationResult, error) {
	var table table.Table

	//is this even valid YAML?
	err := yaml.Unmarshal(yamlBytes, &table)
	if err != nil {
		return nil, err
	}

	//validate the table and parse portions of it since we are tearing the table
	//apart to do the validation anyway
	validationResults := table.Validate()

	//store the table in the repo
	fullName := util.BuildFullName(table.Definition.Namespace, table.Definition.Name, "")
	cr.tableData[fullName] = &table

	return validationResults, nil
}
