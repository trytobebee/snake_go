package main

import (
	"math/rand"
	"time"
)

// NewGame 创建新游戏
func NewGame() *Game {
	g := &Game{
		snake:         []Point{{x: width / 2, y: height / 2}},
		direction:     Point{x: 1, y: 0},
		score:         0,
		gameOver:      false,
		startTime:     time.Now(),
		foodEaten:     0,
		foods:         make([]Food, 0),
		lastFoodSpawn: time.Now(),
	}
	// 初始生成一个豆子
	g.spawnOneFood()
	return g
}

// spawnOneFood 生成一个随机类型的豆子
func (g *Game) spawnOneFood() {
	if len(g.foods) >= maxFoodsOnBoard {
		return // 达到最大数量，不再生成
	}

	// 随机选择豆子类型，高分豆子的概率更低
	randNum := rand.Intn(100)
	var foodType FoodType
	if randNum < 15 { // 15% 概率红色高分
		foodType = FoodRed
	} else if randNum < 35 { // 20% 概率橙色
		foodType = FoodOrange
	} else if randNum < 65 { // 25% 概率蓝色
		foodType = FoodBlue
	} else { // 35% 概率紫色低分
		foodType = FoodPurple
	}

	// 找一个不与蛇身和其他豆子重叠的位置
	for attempts := 0; attempts < 100; attempts++ {
		pos := Point{
			x: rand.Intn(width-2) + 1,
			y: rand.Intn(height-2) + 1,
		}

		// 检查是否与蛇身重叠
		onSnake := false
		for _, p := range g.snake {
			if p == pos {
				onSnake = true
				break
			}
		}
		if onSnake {
			continue
		}

		// 检查是否与其他豆子重叠
		onFood := false
		for _, f := range g.foods {
			if f.pos == pos {
				onFood = true
				break
			}
		}
		if onFood {
			continue
		}

		// 找到合适的位置，生成豆子
		g.foods = append(g.foods, Food{
			pos:       pos,
			foodType:  foodType,
			spawnTime: time.Now(),
		})
		g.lastFoodSpawn = time.Now()
		return
	}
}

// removeExpiredFoods 移除过期的豆子
func (g *Game) removeExpiredFoods() {
	newFoods := make([]Food, 0)
	for _, food := range g.foods {
		if !food.isExpired() {
			newFoods = append(newFoods, food)
		}
	}
	g.foods = newFoods
}

// trySpawnFood 尝试生成新豆子
func (g *Game) trySpawnFood() {
	// 游戏结束时不再生成豆子
	if g.gameOver {
		return
	}

	// 移除过期的豆子
	g.removeExpiredFoods()

	// 如果没有豆子了，立即生成一个
	if len(g.foods) == 0 {
		g.spawnOneFood()
		return
	}

	// 如果距离上次生成已经超过固定周期，且未达到最大数量，生成新豆子
	if time.Since(g.lastFoodSpawn) > foodSpawnInterval && len(g.foods) < maxFoodsOnBoard {
		g.spawnOneFood()
	}
}

// update 更新游戏状态
func (g *Game) update() {
	if g.gameOver {
		return
	}

	// Calculate new head position
	head := g.snake[0]
	newHead := Point{
		x: head.x + g.direction.x,
		y: head.y + g.direction.y,
	}

	// Check wall collision
	if newHead.x <= 0 || newHead.x >= width-1 || newHead.y <= 0 || newHead.y >= height-1 {
		g.gameOver = true
		g.crashPoint = newHead
		return
	}

	// Check self collision
	for _, p := range g.snake {
		if p == newHead {
			g.gameOver = true
			g.crashPoint = newHead
			return
		}
	}

	// Move snake
	g.snake = append([]Point{newHead}, g.snake...)

	// Check food collision - 检查是否吃到任何豆子
	ateFood := false
	for i, food := range g.foods {
		if newHead == food.pos {
			g.score += food.getScore()
			g.foodEaten++ // 增加吃豆计数
			// 移除被吃掉的豆子
			g.foods = append(g.foods[:i], g.foods[i+1:]...)
			ateFood = true
			break
		}
	}

	if !ateFood {
		// Remove tail if no food eaten
		g.snake = g.snake[:len(g.snake)-1]
	}

	// 尝试生成新豆子
	g.trySpawnFood()
}
