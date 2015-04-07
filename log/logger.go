package logger

import (
	l4g "github.com/alecthomas/log4go"
)

var Log = make(l4g.Logger)

func Init() {
	// var config = `{"filename":"server.log","maxdays":30,"daily":true,"rotate":true}`
	// // level:Trace 1|Debug 2|Info 3|Warn 4|Error 5|Critical 6
	// Log := log.NewLogger(10)
	// Log.SetLogger("file", config)
	// Log.EnableFuncCallDepth(true)
	Log.LoadConfiguration("./log4go.xml")
}
