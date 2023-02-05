package console

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var CURRENT_ID int

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {

	var Videos []Video
	var Input_Video Video
	var Author User
	author_id := 0
	/* 	var videoPath string = "files/video_list.json"
	   	video_list_json, err := ioutil.ReadFile(videoPath)
	   	if err != nil {
	   		panic(err)
	   	}
	   	json.Unmarshal(video_list_json, &Videos) */
	rows, err := DB.Query("SELECT * FROM video")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&Input_Video.Id, &Input_Video.PlayUrl, &Input_Video.CoverUrl, &Input_Video.FavoriteCount, &author_id)
		if err != nil {
			panic(err)
		}
		Input_Video.CommentCount = 0
		Input_Video.IsFavorite = false
		user_row, err := DB.Query("SELECT id,name,followcount,followercount FROM user where id=" + strconv.Itoa(author_id))
		if err != nil {
			panic(err)
		}
		for user_row.Next() {
			err = user_row.Scan(&Author.Id, &Author.Name, &Author.FollowCount, &Author.FollowerCount)
			if err != nil {
				panic(err)
			}
			Author.IsFollow = false
		}
		Input_Video.Author = Author
		Videos = append(Videos, Input_Video)
	}
	CURRENT_ID = int(Videos[0].Id)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: Videos,
		NextTime:  time.Now().Unix(),
	})
}
