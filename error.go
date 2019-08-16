package gomysql

import "errors"

var CONNECTDBERROR = errors.New("can't connect db")
var NotConnetKey = errors.New("can't found connect key")
var NotInitERROR = errors.New("not init conf map")
var TAGERROR = errors.New("not found tag")