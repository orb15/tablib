package tablib

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

type luaModule struct {
	repo TableRepository
}

func newLuaModule(r TableRepository) *luaModule {
	return &luaModule{
		repo: r,
	}
}

//LuaModuleLoader registers functions with the lua runtime
func (lm *luaModule) luaModuleLoader(L *lua.LState) int {

	//this map defines functions visible to lua. the key is the name of the
	//function as it is exposed to lua and the value is a pointer to an LGFunction
	//(a function type specified in the gopher-lua lib)
	exportedGoFuncs := map[string]lua.LGFunction{
		"roll": lm.rollOnTable,
	}

	//make certain functions available to lua
	mod := L.SetFuncs(L.NewTable(), exportedGoFuncs)

	// returns the module
	L.Push(mod)
	return 1
}

//RollOnTable is the lua-visible wrapper function for TableRepository.Roll()
func (lm *luaModule) rollOnTable(lState *lua.LState) int {

	//confirm arg is a single string and convert it to a Go string
	argCount := lState.GetTop() //gets count of args passed onto stack
	if argCount != 1 {
		msg := fmt.Sprintf("ERROR: RollOnTable requires 1 argument, received: %d", argCount)
		lState.Push(lua.LString(msg))
		return 1
	}
	tblNameInLuaFmt := lState.Get(1) //lua uses 1-based arrays - get first argument
	tblNameLuaType := tblNameInLuaFmt.Type()
	if tblNameLuaType != lua.LTString {
		msg := fmt.Sprintf("ERROR: RollOnTable requires string argument, received type: %s", tblNameLuaType)
		lState.Push(lua.LString(msg))
		return 1
	}
	tblName := lState.ToString(1)

	//Actually roll on the table specified in the lua script
	tr := lm.repo.Roll(tblName, 1) //always roll once in scripts

	//push the result of the roll back to lua
	lState.Push(lua.LString(tr.Result[0]))
	return 1
}