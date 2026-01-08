package main

import "time"

// getScore 获取豆子的分值
func (f *Food) getScore() int {
	switch f.foodType {
	case FoodRed:
		return 40
	case FoodOrange:
		return 30
	case FoodBlue:
		return 20
	case FoodPurple:
		return 10
	default:
		return 10
	}
}

// getDuration 获取豆子的留存时间
func (f *Food) getDuration() time.Duration {
	switch f.foodType {
	case FoodRed:
		return 10 * time.Second
	case FoodOrange:
		return 15 * time.Second
	case FoodBlue:
		return 18 * time.Second
	case FoodPurple:
		return 20 * time.Second
	default:
		return 20 * time.Second
	}
}

// isExpired 检查豆子是否过期
func (f *Food) isExpired() bool {
	return time.Since(f.spawnTime) > f.getDuration()
}

// getRemainingSeconds 获取剩余秒数
func (f *Food) getRemainingSeconds() int {
	remaining := f.getDuration() - time.Since(f.spawnTime)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// getEmoji 获取豆子的 emoji
func (f *Food) getEmoji() string {
	switch f.foodType {
	case FoodRed:
		return "🔴"
	case FoodOrange:
		return "🟠"
	case FoodBlue:
		return "🔵"
	case FoodPurple:
		return "🟣"
	default:
		return "🟣"
	}
}

// getEmojiWithTimer 获取豆子的 emoji（不带倒计时）
func (f *Food) getEmojiWithTimer() string {
	// 直接返回原始豆子 emoji，倒计时数字将在旁边格子显示
	return f.getEmoji()
}

// getTimerEmoji 获取倒计时数字 emoji（如果在倒计时阶段）
func (f *Food) getTimerEmoji() string {
	remaining := f.getRemainingSeconds()

	// 最后5秒内返回倒计时数字
	// 使用圆圈数字字符，占用全角宽度（2字节）
	if remaining <= 5 && remaining > 0 {
		// 使用圆圈数字，这些是全角字符
		circledNums := map[int]string{
			1: "①",
			2: "②",
			3: "③",
			4: "④",
			5: "⑤",
		}
		return circledNums[remaining]
	}

	return "" // 不在倒计时阶段返回空字符串
}
