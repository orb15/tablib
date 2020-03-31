package tablib

import (
	"tablib/table"
	"tablib/util"

	yaml "gopkg.in/yaml.v2"
)

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
	fullName := util.BuildFullName(table.Definition.Name, "")
	cr.tableData[fullName] = &table

	//if the table has any line tables, add these as well
	if len(table.Inline) > 0 {
		inlines := extractInlineTables(&table)
		for _, ilt := range inlines {
			cr.tableData[ilt.Definition.Name] = ilt
		}
	}

	return validationResults, nil
}

func (cr *concreteTableRepo) List(searchExpr string) ([]*ListResponse, error) {
	return nil, nil
}

func extractInlineTables(mainTable *table.Table) []*table.Table {
	inlinesAsTables := make([]*table.Table, len(mainTable.Inline))
	for _, ilt := range mainTable.Inline {

		def := &table.DefinitionPart{
			Name:      ilt.FullyQualifiedName,
			TableType: "flat",
		}

		content := make([]string, len(ilt.Content))
		for _, c := range ilt.Content {
			content = append(content, c)
		}

		tbl := &table.Table{
			Definition: def,
			RawContent: content,
			IsValid:    true,
		}
		inlinesAsTables = append(inlinesAsTables, tbl)
	}
	return inlinesAsTables
}
