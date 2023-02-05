package main

import (
	"m/console"

	"github.com/gin-gonic/gin"
)

func initApi(r *gin.Engine) {
	r.Static("/static", "./public")
	apiPath := r.Group("/douyin")

	//basic foundations
	apiPath.GET("/feed/", console.Feed)
	apiPath.GET("/user/", console.UserInfo)
	apiPath.POST("/user/register/", console.Register)
	apiPath.POST("/user/login/", console.Login)
	apiPath.POST("/publish/action/", console.Publish)
	apiPath.GET("/publish/list/", console.PublishList)

	//extended foundations - interact
	apiPath.POST("/favorite/action/", console.FavoriteAction)
	apiPath.GET("/favorite/list/", console.FavoriteList)
	apiPath.POST("/comment/action/", console.CommentAction)
	apiPath.GET("/comment/list/", console.CommentList)

	//extended foundations - social
	apiPath.POST("/relation/action/", console.RelationAction)
	apiPath.GET("/relation/follow/list/", console.RelationFollowList)
	apiPath.GET("/relation/follower/list/", console.RelationFollowerList)
	apiPath.GET("/relation/friend/list/", console.RelationFriendList)
	apiPath.POST("/message/action/", console.MessageAction)
	apiPath.GET("/message/chat/", console.MessageChat)
}

func main() {
	console.InitDB()
	// fmt.Printf("console.DB: %v\n", console.DB)
	r := gin.Default()
	initApi(r)
	r.Run()
}
