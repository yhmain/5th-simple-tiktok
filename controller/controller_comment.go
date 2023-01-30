package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"github.com/yhmain/5th-simple-tiktok/util"
)

// 评论操作的返回结构体
type CommentResponse struct {
	util.Response
	Comment model.Comment
}

// 评论列表返回的结构体
type CommentListResponse struct {
	util.Response
	CommentList []model.Comment `json:"comment_list"`
}

//评论操作，路由
func CommentAction(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	video_id := middleware.GetParamPostOrGet(c, "video_id")
	vID, _ := strconv.ParseInt(video_id, 10, 64) // 字符串转化成int64
	action_type := c.Query("action_type")
	if action_type == "1" {
		// 发布评论
		comment_text := middleware.GetParamPostOrGet(c, "comment_text")
		cID := util.GenID() // 生成评论ID
		newComment := model.Comment{
			Id:          cID,
			Content:     comment_text,
			CreatedTime: time.Now().Unix(),
			UserID:      usertoken.UserID,
			VideoID:     vID,
		}
		// 知道了视频的ID，先去查看redis里面是否存了它的评论数，若无则去查mysql数据库
		if _, err := middleware.GetKey("ComCnt:" + video_id); err != nil {
			video := dao.GetVideoByID(vID)                                                   // 去数据库查询 视频结构体
			middleware.SetKey("ComCnt:"+video_id, strconv.FormatInt(video.CommentCount, 10)) // 评论数量加入redis
		}
		// 记录插入redis
		middleware.UpdateRedisComment(action_type, newComment)
		// 返回
		c.JSON(http.StatusOK, CommentResponse{
			Response: util.Success, //成功
			Comment:  newComment,
		})
	} else {
		//删除评论
		comment_id := middleware.GetParamPostOrGet(c, "comment_id")
		cid, _ := strconv.ParseInt(comment_id, 10, 64)
		newComment := model.Comment{
			Id:      cid,
			UserID:  usertoken.UserID,
			VideoID: vID,
		}
		// 知道了视频的ID，先去查看redis里面是否存了它的评论数，若无则去查mysql数据库
		if _, err := middleware.GetKey("ComCnt:" + video_id); err != nil {
			video := dao.GetVideoByID(vID)                                                   // 去数据库查询 视频结构体
			middleware.SetKey("ComCnt:"+video_id, strconv.FormatInt(video.CommentCount, 10)) // 评论数量加入redis
		}
		// 记录更新到redis
		middleware.UpdateRedisComment(action_type, newComment)
		// 返回
		c.JSON(http.StatusOK, CommentResponse{
			Response: util.Success, //成功
			Comment:  model.Comment{},
		})
	}

}

//评论列表，路由
func CommentList(c *gin.Context) {
	// 查询之前，将Redis数据手段更新到Mysql
	SaveRedisToMySQL()
	// 再查询MySQL数据
	video_id := middleware.GetParamPostOrGet(c, "video_id")
	vID, _ := strconv.ParseInt(video_id, 10, 64) // 字符串转化成int64
	comments, err := dao.GetCommentsByVid(vID)   // 查询mysql数据
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    util.ComListErr,
			CommentList: []model.Comment{},
		})
		return
	}
	// 返回
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    util.Success,
		CommentList: comments,
	})
}
