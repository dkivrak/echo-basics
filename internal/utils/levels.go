package utils

func GetLogLevel(flag FlagEnum) int {
	switch flag {
	case LogFlag:
		return 0
	case DebugFlag:
		return 1
	case InfoFlag:
		return 2
	case WarnFlag:
		return 3
	case ErrorFlag:
		return 4
	case TraceFlag:
		return 5
	default:
		return -1
	}
}
