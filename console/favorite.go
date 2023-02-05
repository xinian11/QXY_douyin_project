package console

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	var Temp_user User
	token := c.Query("token")
	if token == "" {
		panic("没读到")
	}
	video_id := c.Query("video_id")
	if video_id == "" {
		panic("没读到")
	}
	action_type := c.Query("action_type")
	if action_type == "" {
		panic("没读到")
	}
	Temp_user.Id = 0
	exist, err := DB.Query("select id,name,followcount,followercount from user where token=\"" + token + "\"")
	if err != nil {
		panic(err)
	}
	for exist.Next() {
		err = exist.Scan(&Temp_user.Id, &Temp_user.Name, &Temp_user.FollowCount, &Temp_user.FollowerCount)
		if err != nil {
			panic(err)
		}
		Temp_user.IsFollow = false
		Current_user.User = Temp_user

	}
	if Temp_user.Id == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	favorite_change, err := DB.Begin()
	if err != nil {
		panic(err)
	}
	if action_type == "1" {
		r, err := DB.Exec("update video set favoritecount=favoritecount+1 where id=" + video_id)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			favorite_change.Rollback()
			return
		}
		r, err = DB.Exec("insert into favorite(user_id,video_id)values(?,?)", Temp_user.Id, video_id)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			favorite_change.Rollback()
			return
		}
		favorite_change.Commit()
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "点赞成功",
		})
	}
	if action_type == "2" {
		r, err := DB.Exec("update video set favoritecount=favoritecount-1 where id=" + video_id)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			favorite_change.Rollback()
			return
		}
		r, err = DB.Exec("delete from favorite where video_id=" + video_id + " and user_id=" + strconv.Itoa(int(Temp_user.Id)))
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			favorite_change.Rollback()
			return
		}
		favorite_change.Commit()
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取消点赞",
		})
	}
}

func FavoriteList(c *gin.Context) {
	var Current_User_Favorite_Video []Video
	var Temp_video Video
	var Favorite_video_Id []int
	var Temp_id int
	user_id := c.Query("user_id")
	rows, err := DB.Query("select video_id from favorite where user_id=" + user_id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&Temp_id)
		if err != nil {
			panic(err)
		}
		Favorite_video_Id = append(Favorite_video_Id, Temp_id)
	}
	// fmt.Printf("Favorite_video_Id[len(Favorite_video_Id)-1]: %v\n", Favorite_video_Id[len(Favorite_video_Id)-1])
	for _, i := range Favorite_video_Id {
		rows, err = DB.Query("select playurl,coverurl,favoritecount from video where id=" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			err = rows.Scan(&Temp_video.PlayUrl, &Temp_video.CoverUrl, &Temp_video.FavoriteCount)
			if err != nil {
				panic(err)
			}
			Temp_video.Id = int64(i)
			Temp_video.Author = Current_user.User
			Temp_video.CommentCount = 0
			Temp_video.IsFavorite = true
			Current_User_Favorite_Video = append(Current_User_Favorite_Video, Temp_video)
		}
	}
	// fmt.Printf("Current_User_Favorite_Video[len(Current_User_Favorite_Video)-1]: %v\n", Current_User_Favorite_Video[len(Current_User_Favorite_Video)-1])
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: Current_User_Favorite_Video,
	})
}
