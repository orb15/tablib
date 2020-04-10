package tablib

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func executeScript(scriptName string, nameSvc nameResolver, repo TableRepository) map[string]string {
	//set up a new lua VM
	//TODO: limit call stack and repository sizes, prevent use of lua modules
	//that allow access to OS, filesys or other dangerous crap
	lState := lua.NewState()
	defer lState.Close()

	//tell the lua VM about the go code we are exposing to it
	luaMod := newLuaModule(repo)
	lState.PreloadModule("tables", luaMod.luaModuleLoader)

	//fetch the precompiled lua script by name
	scriptData, err := nameSvc.scriptForName(scriptName)
	if err != nil {
		createErrorMap(scriptName, fmt.Sprintf("%s", err))
	}

	//prep the lua script - all this does is store the precompiled code
	//in lState and await our call to lua functions it defines
	luafunc := lState.NewFunctionFromProto(scriptData)
	lState.Push(luafunc)
	err = lState.PCall(0, lua.MultRet, nil)
	if err != nil {
		return createErrorMap(scriptName,
			fmt.Sprintf("fail to load loading compiled script: %s", err))
	}

	//TODO: here we need to call well-known lua function to get info about
	//the params the lua main() code needs to do its job. Once we get these, this
	//method (Execute) will need to utilize a callback function (needs passed in)
	//to request the param values from the caller of this lib.

	//For sanity sake, all lua functions should take and return a single well-known
	//type so we always know the size of the argument list being passed or
	//returned. A map[string]string is sufficent and simple to handle

	ldm := make(map[string]string) //hack: make up params to pass for now
	ldm["p1"] = "v1"
	ldm["p2"] = "v2"

	//call the well-known function "main" which is the 'main' for our lua script
	if err := lState.CallByParam(lua.P{
		Fn:      lState.GetGlobal("main"),
		NRet:    0,
		Protect: true,
	}, toLuaLTable(ldm)); err != nil {
		if err != nil {
			return createErrorMap(scriptName, fmt.Sprintf("executing main(): %s", err))
		}
	}

	//retrieve the well-known return value from lua
	retval := lState.GetGlobal("rettbl")
	retmap := fromLuaTable(scriptName, lState, retval)

	return retmap
}

//converts a go map to a lua LTable
func toLuaLTable(goMap map[string]string) *lua.LTable {
	ltbl := &lua.LTable{}
	for k, v := range goMap {
		ltbl.RawSetString(k, lua.LString(v))
	}
	return ltbl
}

//converts a lua LTable to a go map
func fromLuaTable(scriptName string, lState *lua.LState, lVal lua.LValue) map[string]string {

	//do we really have an LTable in the passed LValue?
	if lVal.Type() != lua.LTTable {
		return createErrorMap(scriptName,
			"script does not contain required return table variable 'rettbl'")
	}
	luaTable := lVal.(*lua.LTable)

	mp := make(map[string]string)

	luaTable.ForEach(func(k lua.LValue, v lua.LValue) {
		key := k.String()
		val := v.String()
		mp[key] = val
	})
	return mp
}

//helper to uniformly return errors during script execution
func createErrorMap(scriptName, details string) map[string]string {
	errMap := make(map[string]string)
	errMap["Script-Error"] = details
	return errMap
}
