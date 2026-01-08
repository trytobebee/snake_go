package main

import "time"

// Point 坐标点
type Point struct {
	x, y int
}

// FoodType 豆子类型
type FoodType int

const (
	FoodPurple FoodType = iota // 紫色，10分，20s
	FoodBlue                   // 蓝色，20分，18s
	FoodOrange                 // 橙色，30分，15s
	FoodRed                    // 红色，40分，10s
)

// Food 豆子结构
type Food struct {
	pos       Point
	foodType  FoodType
	spawnTime time.Time
}

// Game 游戏主结构
type Game struct {
	snake         []Point
	foods         []Food // 多个豆子
	direction     Point
	score         int
	gameOver      bool
	paused        bool          // 暂停状态
	crashPoint    Point         // 碰撞位置
	startTime     time.Time     // 游戏开始时间
	foodEaten     int           // 吃到的豆子数量
	pausedTime    time.Duration // 暂停累计时间
	pauseStart    time.Time     // 暂停开始时间
	lastFoodSpawn time.Time     // 上次生成豆子的时间
}
