package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
)

func TestGetRedisClient(t *testing.T) { // 测试ping，获取redis连接
	_, err := GetRedisClient()
	var expectedResult error
	assert.Equal(t, expectedResult, err)
}

func TestGetRedisLike(t *testing.T) {
	UpdateRedisLike("Fav:2:1618519735576563712", "1618519735576563712", "1")
	likes, err := GetRedisLike()
	assert.Equal(t, nil, err)
	expectedResult := []model.Like{
		{
			Id:         "2:1618519735576563712",
			UserID:     int64(2),
			VideoID:    int64(1618519735576563712),
			IsFavorite: true,
		},
	}
	assert.Equal(t, expectedResult, likes)
}
