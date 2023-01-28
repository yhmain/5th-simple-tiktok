package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/controller"
	"github.com/yhmain/5th-simple-tiktok/middleware"
)

func initRouter(r *gin.Engine) {

	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.GET("/user/", middleware.JWTAuth(), controller.UserInfo)
	apiRouter.POST("/publish/action/", middleware.JWTAuth(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.JWTAuth(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.JWTAuth(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.JWTAuth(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.JWTAuth(), controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// // extra apis - II
	// apiRouter.POST("/relation/action/", controller.RelationAction)
	// apiRouter.GET("/relation/follow/list/", controller.FollowList)
	// apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	// apiRouter.GET("/relation/friend/list/", controller.FriendList)
	// apiRouter.GET("/message/chat/", controller.MessageChat)
	// apiRouter.POST("/message/action/", controller.MessageAction)
}
