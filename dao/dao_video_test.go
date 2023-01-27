package dao

import (
	"fmt"
	"testing"

	"github.com/yhmain/5th-simple-tiktok/model"
)

func TestInsertNewVideo(t *testing.T) {
	video := model.Video{
		Id:      90349901223,
		Title:   "lalal",
		PlayUrl: "playurl.mp4",
		UserID:  2,
	}
	InsertNewVideo(&video)
}

func TestGetVideoByIDs(t *testing.T) {
	vs := GetVideoByIDs([]int64{int64(1618519735576563712), int64(1618885198131236864)})
	fmt.Println(vs)
}
