package dao

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
)

func TestInsertNewComment(t *testing.T) {
	comment := &model.Comment{
		Id:          123456789,
		Content:     "测试的自动插入dd",
		CreatedTime: 1674807386,
		UserID:      2,
		VideoID:     1618885198131236864,
	}
	err := InsertNewComment(comment)
	fmt.Println(err)
	assert.Equal(t, nil, err)
}

func TestGetCommentsByVid(t *testing.T) {
	cos, err := GetCommentsByVid(int64(1618885198131236864))
	fmt.Println(cos)
	assert.Equal(t, nil, err)
}
