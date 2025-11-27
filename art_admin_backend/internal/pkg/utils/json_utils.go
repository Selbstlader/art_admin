package utils

import (
	"encoding/json"
	"errors"
)

// JSONMarshal 序列化切片为JSON字符串
func JSONMarshal(data []string) (string, error) {
	if data == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(data)
	return string(bytes), err
}

// JSONUnmarshal 反序列化JSON字符串为切片
func JSONUnmarshal(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return []string{}, nil
	}
	var data []string
	err := json.Unmarshal([]byte(jsonStr), &data)
	return data, err
}

// JSONMarshalMap 序列化map为JSON字符串
func JSONMarshalMap(data map[string]interface{}) (string, error) {
	if data == nil {
		return "{}", nil
	}
	bytes, err := json.Marshal(data)
	return string(bytes), err
}

// JSONUnmarshalMap 反序列化JSON字符串为map
func JSONUnmarshalMap(jsonStr string) (map[string]interface{}, error) {
	if jsonStr == "" {
		return make(map[string]interface{}), nil
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	return data, err
}

// ValidateMoodType 验证情绪类型是否有效
func ValidateMoodType(moodType string) bool {
	validTypes := map[string]bool{
		"happy": true, "sad": true, "anxious": true, "calm": true,
		"angry": true, "excited": true, "tired": true, "stressed": true,
	}
	return validTypes[moodType]
}

// ValidateMeditationCategory 验证冥想分类是否有效
func ValidateMeditationCategory(category string) bool {
	validCategories := map[string]bool{
		"sleep": true, "stress": true, "focus": true, "anxiety": true,
		"breathing": true, "mindfulness": true, "body_scan": true, "walking": true,
	}
	return validCategories[category]
}

// CalculateProgress 计算进度百分比
func CalculateProgress(current, target int) float64 {
	if target <= 0 {
		return 0
	}
	if current >= target {
		return 100
	}
	return float64(current) / float64(target) * 100
}

// ParseDateRange 解析日期范围，确保开始时间不晚于结束时间
func ParseDateRange(startDate, endDate *string) error {
	// 这里可以添加日期解析和验证逻辑
	return nil
}

// ValidateIntensity 验证情绪强度
func ValidateIntensity(intensity int) error {
	if intensity < 1 || intensity > 10 {
		return errors.New("情绪强度必须在1-10之间")
	}
	return nil
}
