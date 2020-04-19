package tablib

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"tablib/tableresult"
	"tablib/validate"
)

//TableRepository holds the various tables and Lua scripts and manages their execution
type TableRepository interface {

	//AddLuaScript stores the given lua script in the repository.
	//The script name and the script string itself are mandatory.
	AddLuaScript(scriptName string, luaScript string) error

	//AddTable stores the given yaml representation of a table in the repository.
	//
	//If the presnted yaml is not parsable or has other structural issues, an error is raised.
	//Errors and warnings related to the semantics of the table (e.g. internal consistency
	//issues or table syntax errors) are captured in the returned ValidationResult
	AddTable(yamlBytes []byte) (*validate.ValidationResult, error)

	//Execute executes the named Lua script and returns a map of named output keys and their values.
	//
	//The nature of the output is arbitrary and is defined by the Lua script, but typically
	//the returned map contains name/value pairs where the name represents some meaningful
	//table roll and the value is the results from that table roll. The map may contain any other
	//information the script wishes to communicate to the client, like a synthesized or concatenated
	//result of several tables rolls that have been assembled in the script to produce an overall
	//meaningful result. If an error occurs in script execution, error information will be returned
	//in the map in place of the intended result(s).
	//
	//The callback function is optional. If used, it will be called if the Lua script requests
	//parameters from the caller. If set to nil, the Lua script will be returned the default
	//value of each parameter. For more information on the format of Lua scripts and the
	//use of the callback, see the README documentation.
	Execute(scriptName string, callback ParamSpecificationRequestCallback) map[string]string

	//EvaluateDiceExpression revaluates a dice expression and returns the result or
	//an a non-nuil ValidatioNresult if the expression is not valid.
	EvaluateDiceExpression(diceExpr string) (int, error)

	//List provides the raw string listing of the named table or script.
	//
	//An error is returned if the named item does not exist or if itemType is anything
	//other than "table" or "script"
	List(name string, itemType string) (string, error)

	//Pick returns count unique items from the named table.

	//The table type must be flat; providing the name of a ranged table will generate
	//an error
	Pick(tableName string, count int) *tableresult.TableResult

	//Roll 'rolls' on the named table count times, generating a single result with each roll
	Roll(tableName string, count int) *tableresult.TableResult

	//Search returns information about the tables and scripts in the repository.

	//The namePredicate must be a valid regular expression and is optional. If not provided,
	//an item's name will not be used as filtering criteria. In a similar fashion, tags
	//is also optional and if not provided, will not be used as filtering criteria.

	//To obtain the entire contents of the repository, call Search("", nil).
	//To filter by tags only, call Search ("", []string{"tag1", "tag2"}). This will
	//return all items in the repo that have at least one of these tags defined.
	//To filter by name only, call Search("myregex", nil). This will return all items
	//that match the given regex.  Provide both parameters for a narrow search. Note
	//that in this case, the results are filtered first by tag and then by name.
	Search(namePredicate string, tags []string) ([]*SearchResult, error)
}

//SearchResult holds information about each object discovered during a search
type SearchResult struct {
	Name string
	Type string
	Tags []string
}

func (sr *SearchResult) toFullComparable() string {
	sort.Strings(sr.Tags)
	return fmt.Sprintf("%s:%s:%s", sr.Name, sr.Type, strings.Join(sr.Tags, ":"))
}
func (sr *SearchResult) toComparable() string {
	return fmt.Sprintf("%s:%s", sr.Name, sr.Type)
}

//NewTableRepository does what it says on the tin
func NewTableRepository() TableRepository {
	return &concreteTableRepo{
		tableStore:      make(map[string]*tableData),
		scriptStore:     make(map[string]*scriptData),
		tagSearchCache:  make(map[string][]*SearchResult),
		nameSearchCache: make(map[string]*SearchResult),
		lock:            &sync.RWMutex{},
	}
}
