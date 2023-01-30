package controller

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/config"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"github.com/yhmain/5th-simple-tiktok/util"
)

//  /feed接口的响应体
type FeedResponse struct {
	util.Response               //标准响应体
	VideoList     []model.Video `json:"video_list"`          //视频列表
	NextTime      int64         `json:"next_time,omitempty"` //本次返回的视频中，发布最早的时间，作为下次请求的Latest_time
}

// 视频列表的响应体
type VideoListResponse struct {
	util.Response
	VideoList []model.Video `json:"video_list"`
}

// 如果当前是登录状态，则更新点赞信息（是否点赞和点赞数量）、评论信息（评论内容和评论数量）
// 还有是因为要取得redis最新的数据，点赞数量和评论数量
func updateVideoInfo(userID int64, videos []model.Video) {
	for i := 0; i < len(videos); i++ {
		// 更新点赞状态，查询redis
		likeID := fmt.Sprintf("%d:%d", userID, videos[i].Id)
		if val, err := middleware.GetKey(likeID); err != nil { //  若redis里面有数据，则以之为准
			if val == "1" {
				videos[i].IsFavorite = true //如果为1表示点赞了
			}
		}
		// 更新视频的点赞数量，查询redis计数器
		video_id := strconv.FormatInt(videos[i].Id, 10)
		vidKey := "FavCnt:" + video_id                         // 构造redis key
		if val, err := middleware.GetKey(vidKey); err != nil { //  若redis里面有数据，则以之为准
			cnt, _ := strconv.ParseInt(val, 10, 64)
			videos[i].FavoriteCount = cnt // 最终更新到 videos中
		}
		// 更新视频的评论数量，查询redis计数器
		comCntKey := "ComCnt:" + video_id
		if val, err := middleware.GetKey(comCntKey); err != nil { //  若redis里面有数据，则以之为准
			cnt, _ := strconv.ParseInt(val, 10, 64)
			videos[i].CommentCount = cnt // 最终更新到 videos中
		}
	}
}

//视频流接口，路由
func Feed(c *gin.Context) {
	t := c.Query("latest_time")
	latest_time, err := strconv.ParseInt(t, 10, 64) //string转化为int64
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  util.InvalidTimeErr, //失败
			VideoList: nil,
			NextTime:  time.Now().Unix(),
		})
	} else {
		if t == "" { //空字符串会转化为0，则表示取当前时间的时间戳
			latest_time = time.Now().Unix()
		}
		//获取视频数据
		var videos = dao.GetVideosByTime(latest_time)
		var nextTime = time.Now().Unix()
		if len(videos) >= 1 { // 注意：返回的视频列表可能为空
			nextTime = videos[len(videos)-1].CreatedTime
			token := c.Query("token") // 如果当前是用户登录状态，则token存在，需要查询 当前用户是否喜欢该视频
			if token != "" {
				_, claims, err := middleware.ParseToken(token)
				if err != nil {
					fmt.Println("Token解析出错: ", err, "出错的Token是Feed里面：", token)
					c.JSON(http.StatusOK, FeedResponse{
						Response:  util.ParseTokenErr, //token解析失败
						VideoList: []model.Video{},
						NextTime:  time.Now().Unix(), // 当前时间
					})
					return
				}
				uid := claims.UserID // 得到了当前用户ID
				// fmt.Println(videos)
				updateVideoInfo(uid, videos) // 更新点赞信息和评论数量信息
			} else {
				// 如果没有用户登录，则点赞状态都为false
				for i := 0; i < len(videos); i++ {
					videos[i].IsFavorite = false //首先默认是false
					// fmt.Println(videos[i])
				}
			}

		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  util.Success, //成功
			VideoList: videos,
			NextTime:  nextTime, //本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
		})
	}
}

// 投稿接口，用户发布视频，路由
func Publish(c *gin.Context) {
	usertoken := c.MustGet("usertoken").(middleware.UserToken) //经过jwt鉴权后解析出的usertoekn
	fmt.Println("111")
	// 获取视频流数据
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, util.FileParseErr) //文件解析失败
		fmt.Println("222")
		return
	}
	fileExt := path.Ext(file.Filename)       // 获取文件后缀名，比如 .mp4  .jpg
	fileName := util.GenFileName() + fileExt // 利用时间戳构造新的文件名
	//构建Video结构体所需要的参数
	newVideoID := util.GenID()        //获取新视频的ID
	paramTitle := c.PostForm("title") //视频标题
	playUrl := fileName               //视频播放路径，因为视频都是放在一个目录下，所以视频名也要确保唯一

	//上传的视频保存到本地服务器
	saveFile := filepath.Join("./public/video/", playUrl)
	// gin 简单做了封装,拷贝了文件流
	if err := c.SaveUploadedFile(file, saveFile); err != nil {
		c.JSON(http.StatusOK, util.VideoUploadErr) //上传视频失败
		return
	}
	//从上传到本地的视频中抽取一帧作为封面，并保存到服务器（本地）
	coverUrl := strconv.FormatInt(newVideoID, 10) + ".jpeg"                    //利用新生成的视频ID， 构造视频封面路径
	saveImg := filepath.Join("./public/img/", coverUrl)                        //调用ffmpeg对应的 图片生成路径
	if msg, resp := util.GetVideoFrame(saveFile, saveImg, 3); msg == "Error" { // 这里设置抽取第3帧作为封面
		c.JSON(http.StatusOK, resp) //截取封面失败
		return
	}
	//向数据库里面插入记录
	createdTime := time.Now().Unix() //获取当前时间戳
	userId := usertoken.UserID       //发布视频的用户ID
	// 构造保存到数据库的前缀播放路径
	params := config.GetConfig().App
	videoPrefix := fmt.Sprintf("http://%s:%s/static/%s/", params.Host, params.Port, params.Video)
	imgPrefix := fmt.Sprintf("http://%s:%s/static/%s/", params.Host, params.Port, params.Img)
	playUrl = videoPrefix + playUrl
	coverUrl = imgPrefix + coverUrl
	// 构造结构体数据
	newVideo := model.Video{Id: newVideoID, Title: paramTitle, PlayUrl: playUrl,
		CoverUrl: coverUrl, CreatedTime: createdTime, UserID: userId}
	if err := dao.InsertNewVideo(&newVideo); err != nil {
		c.JSON(http.StatusOK, util.VideoInsertErr) //上传文件失败
		return

	}
	//上传成功！
	c.JSON(http.StatusOK, util.Success.WithMessage(playUrl+" 视频上传成功！"))
}

// 发布列表的接口，路由
func PublishList(c *gin.Context) {
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	uid := usertoken.UserID

	//调用service 获取该用户发布的视频列表
	var videos = dao.GetVideosByUserID(uid)
	// 更新点赞信息
	updateVideoInfo(uid, videos)

	//返回视频列表和状态码
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  util.Success,
		VideoList: videos,
	})
}
