package logger

import log "github.com/astaxie/beego/logs"

var Log *log.BeeLogger

func Init() {
	var config = `{"filename":"server.log","maxdays":30}`
	// levle:Trace 5|Debug 4|Info 3|Warn 2|Error 1|Critical 0
	Log = log.NewLogger(10000)
	Log.SetLogger("file", config)
	Log.EnableFuncCallDepth(false)
	Log.Info("Logger Started")
}
