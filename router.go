package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/controller"
	"github.com/yhmain/5th-simple-tiktok/middleware"
)

func initRouter(r *gin.Engine) {

	// public directory is used to serve static resources
	r.Static("/static", "./public")
	// r.StaticFS("/static", http.Dir("./public"))

	authRouter := r.Group("/douyin")
	authRouter.Use(middleware.JWTAuth()) // 需要经过JWT鉴权
	{
		// basic apis
		authRouter.GET("/user/", controller.UserInfo)
		authRouter.POST("/publish/action/", controller.Publish)
		authRouter.GET("/publish/list/", controller.PublishList)

		// extra apis - I
		authRouter.POST("/favorite/action/", controller.FavoriteAction)
		authRouter.GET("/favorite/list/", controller.FavoriteList)
		authRouter.POST("/comment/action/", controller.CommentAction)

		// extra apis - II
		authRouter.POST("/relation/action/", controller.RelationAction)
		authRouter.GET("/relation/follow/list/", controller.FollowList)
		authRouter.GET("/relation/follower/list/", controller.FollowerList)
		authRouter.GET("/relation/friend/list/", controller.FriendList)
		authRouter.GET("/message/chat/", controller.MessageChat)
		authRouter.POST("/message/action/", controller.MessageAction)
	}

	apiRouter := r.Group("/douyin")
	{
		// basic apis
		apiRouter.GET("/feed/", controller.Feed)
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)

		// extra apis - I
		apiRouter.GET("/comment/list/", controller.CommentList)

		//  extra apis - II
		// ...
	}

}
