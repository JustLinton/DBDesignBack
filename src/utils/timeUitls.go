package utils

const (
	// SecondsPerMinute 定义每分钟的秒数
	SecondsPerMinute = 60
	// SecondsPerHour 定义每小时的秒数
	SecondsPerHour = SecondsPerMinute * 60
	// SecondsPerDay 定义每天的秒数
	SecondsPerDay = SecondsPerHour * 24
)

// ResolveTime 将传入的“秒”解析为3种时间单位
func ResolveTime(seconds int) (day int, hour int, minute int) {
	day = seconds / SecondsPerDay
	hour = seconds / SecondsPerHour
	minute = seconds / SecondsPerMinute
	return
}