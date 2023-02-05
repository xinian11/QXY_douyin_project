package console

import (
	"fmt"
	"m/util"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

func Publish(c *gin.Context) {
	var Temp_user User
	token := c.PostForm("token")
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
	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		panic(err)
	}
	name := filepath.Base(data.Filename)
	videoName := name + c.GetString("FileType")
	coverName := name + ".jpg"
	videoSavePath := "./public/video/" + videoName
	coverSavePath := "./public/cover/" + coverName
	user := Current_user.User
	if err := c.SaveUploadedFile(data, videoSavePath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		panic(err)
	}
	if err = util.GetFrame(videoSavePath, coverSavePath); err != nil {
		// 封面无法保存
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		panic(err)
	}
	video_insect, err := DB.Begin()
	if err != nil {
		panic(err)
	}
	playUrl := "http://" + c.Request.Host + "/static/video/" + videoName
	coverUrl := "http://" + c.Request.Host + "/static/cover/" + coverName
	r, err := DB.Exec("INSERT INTO video(playurl,coverurl,favoritecount,author_id)VALUES(?,?,?,?)", playUrl, coverUrl, 0, user.Id)
	if err != nil {
		panic(err)
	}
	_, err = r.LastInsertId()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		video_insect.Rollback()
		return
	}
	video_insect.Commit()
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  videoName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var Current_User_Publish_Video []Video
	var Temp_video Video
	rows, err := DB.Query("select id,playurl,coverurl,favoritecount from video where author_id=" + strconv.Itoa(int(Current_user.User.Id)))
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&Temp_video.Id, &Temp_video.PlayUrl, &Temp_video.CoverUrl, &Temp_video.FavoriteCount)
		if err != nil {
			panic(err)
		}
		Temp_video.Author = Current_user.User
		Temp_video.CommentCount = 0
		Temp_video.IsFavorite = false
		Current_User_Publish_Video = append(Current_User_Publish_Video, Temp_video)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: Current_User_Publish_Video,
	})
}
