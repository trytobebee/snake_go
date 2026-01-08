package main

import "time"

const (
	width  = 25 // 减小宽度，因为 emoji 是双倍宽度
	height = 25

	// 豆子生成周期
	foodSpawnInterval = 5 * time.Second
	maxFoodsOnBoard   = 5 // 地图上最多同时存在的豆子数量

	// 加速相关常量
	baseTick             = 50 * time.Millisecond  // 基础 tick 间隔
	normalTicksPerUpdate = 3                      // 正常速度：50ms * 3 = 150ms
	boostTicksPerUpdate  = 1                      // 加速：50ms * 1 = 50ms (3倍速)
	boostTimeout         = 150 * time.Millisecond // 加速超时时间
	boostThreshold       = 2                      // 需要连续按几次才触发加速
	keyRepeatWindow      = 200 * time.Millisecond // 判定连续按键的时间窗口
)
