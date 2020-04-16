# Tablib

Tablib is short for "table library".  It is a library to execute an 'extensible random table generation engine' for use in tabletop gaming. Tablib is written in Go and is intended to serve as the backend for web services, desktop applications or to be embedded in other applications needing random table support.

## Overview
The library offers:
* An open source, easy-to-use [Table](#tab-ref) definition language
* A powerful, robust [API](#api-ref)
* A Lua script execution [engine](#lua-ref) for more sophisticated needs

#### Defining Tables
Tables are defined in [YAML](https://yaml.org/spec/1.2/spec.html#Introduction) and look like this (see the [Table Reference](#tab-ref) for the full table syntax):
```
definition:
  name: Ice_cream_flavors
  type: flat
content:
  - chocolate
  - vanilla
  - strawberry
  - fudge swirl
```

Tables are loaded into a TableRepository and can be located and accessed with the [API](#api-ref):
```
package main
import (
  "fmt"
  "tablib"
  )

func main() {

  //create a new repository. All tables are stored in a repo
  repo := NewTableRepository()

  //load and validate the YAML file. Provides extensive consistency and
  //validity responses not shown here
  repo.AddTable("path/to/icecream.yml")

  //roll twice on the specified table and return the values as a slice.
  //Provides graceful error handling and a Log that describes how the result
  //was achieved          
  tabData := repo.Roll("Ice_cream_flavors", 2)

  for _, r := range tabData.Result {
    fmt.Printf("A random roll on the ice cream table: %s\n", r)
  }
}
```
This code will output something like:
```
A random roll on the ice cream table: chocolate
A random roll on the ice cream table: vanilla
```
Mechanisms are provided for picking a specified number of unique values from a table, for ranged or weighted tables (e.g. tables that do not have an equal distribution of results), for tables calling other tables to retrieve data and for tables to declare 'inline tables' when a table's contents needs to be flexible but the flexibility does not warrant the creation of a new table in its own right.

#### Lua Scripting
Tablib also provides a [Lua](http://www.lua.org/about.html) script execution engine that serves as a powerful tool to generate sophisticated results from the tables and stitch together the results of many table rolls. See the [Lua Reference](#lua-ref) section for more details.
```
local t = require("tables")

results {}
function main()
  results["ice-cream-flavor"] = t.roll("Ice_cream_flavors")
  results["syrup"] = t.roll("sundae-syrup")
  results["toppings"] = t.pick("sundae-toppings", 3)
end
```
Scripts are loaded into a TableRepository and can be located and accessed with the [API](#api-ref):
```
package main
import (
  "fmt"
  "tablib"
  )

func main() {

  //create a new repository. All scripts are stored in a repo
  repo := NewTableRepository()

  //load and validate the YAML file. Provides extensive consistency and validity
  //responses not shown here
  repo.AddLuaScript(luaCodeAsString)

  //execute the script. Provides graceful error handling and many other features          
  tabMap := repo.Execute("sundae", nil)

  for key, value := range tabMap {
    fmt.Printf("%s: %s\n", key, value)
  }
}
```

The output of the Lua script is a map (associative array, dictionary table):
```
ice-cream-flavor: chocolate
syrup: hot fudge
toppings: nuts|whipped cream|sprinkles
```
The Lua engine supports a wide range of Lua functions and a parameter callback mechanism to enable the script to request data from the library's caller.

## Inspiration
The Windows Desktop application [Inspiration Pad Pro3](http://www.nbos.com/products/inspiration-pad-pro) is an excellent
random table execution engine. I have made extensive use of this software to create my own tables (see my [IPP3 project](https://github.com/orb15/ipp3) for the kinds of tables I have created). Over time, I found myself wanting to free myself of the desktop client and access this data over the web. I also wanted to embed random tables in other software I was writing and create sophisticated tables for which the IPP3 product had either no or a very limited syntax. Tablib is my solution to these needs.

## <a name="tab-ref"></a> Table Reference
TODO

## <a name="api-ref"></a> API Reference
TODO

## <a name="lua-ref"></a> Lua Reference
TODO
