package api

import "math"

// 布尔值常量
const (
	True  = 1 // 表示真
	False = 0 // 表示假
)

// TimeoutMaxSeconds 约 24 天
const TimeoutMaxSeconds = math.MaxInt32 / 1000

// 令牌状态常量
const (
	TokenStatusInit     = 60011001 // 初始状态
	TokenStatusQueueing = 60011041 // 排队状态
)
