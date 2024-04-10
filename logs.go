package winterSocket

var logDebug func(string)
var logInfo func(string)
var logError func(string, any)

func SetLog(logDebug_ func(string), logInfo_ func(string), logError_ func(string, any)) {
	logDebug = logDebug_
	logInfo = logInfo_
	logError = logError_
}

func pDebug(msg string) {
	if nil == logDebug {
		return
	}
	logDebug(msg)
}

func pInfo(msg string) {
	if nil == logInfo {
		return
	}
	logInfo(msg)
}

func pError(msg string, e any) {
	if nil == logError {
		return
	}
	logError(msg, e)
}
