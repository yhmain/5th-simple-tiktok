package dao

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
)

func TestUpsetLikes(t *testing.T) {
	likes := []model.Like{
		{
			Id:         "2:1618519735576563712",
			UserID:     int64(2),
			VideoID:    int64(1618519735576563712),
			IsFavorite: true,
		},
	}
	err := SaveLikes(likes, nil)
	fmt.Println(err)
	assert.Equal(t, nil, err)
}

func TestGetLikeByID(t *testing.T) {
	like := GetLikeByID("2:1618519735576563712")
	expectedResult := model.Like{Id: "2:1618519735576563712", UserID: int64(2), VideoID: int64(1618519735576563712), IsFavorite: true}
	assert.Equal(t, expectedResult, like)
}
