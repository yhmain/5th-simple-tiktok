package controller

import (
	"testing"

	"github.com/yhmain/5th-simple-tiktok/middlerware"
)

func TestFavoriteList(t *testing.T) {
	middlerware.UpdateRedisLike("Fav:2:1618519735576563712", "1618519735576563712", "1")
	SaveRedisToMySQL()
}
