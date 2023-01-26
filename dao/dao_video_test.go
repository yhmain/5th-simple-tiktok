package dao

import (
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
