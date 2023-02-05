package console

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
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
	if action_type == "1" {
		comment_text := c.Query("comment_text")
		if len(comment_text) > 100 {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论字数过长。"})
			return
		}
		comment_change, err := DB.Begin()
		if err != nil {
			panic(err)
		}
		time_insect := strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(int(time.Now().Day()))
		r, err := DB.Exec("insert into comment(content,createdate,comment_user_id,video_id)values(?,?,?,?)", comment_text, time_insect, Temp_user.Id, video_id)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			comment_change.Rollback()
			return
		}
		comment_change.Commit()
	}
	if action_type == "2" {
		comment_id := c.Query("comment_id")
		comment_change, err := DB.Begin()
		if err != nil {
			panic(err)
		}
		_, err = DB.Exec("delete from comment where id=" + comment_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "评论删除失败。"})
			comment_change.Rollback()
			return
		}
		comment_change.Commit()
	}
}

func CommentList(c *gin.Context) {
	var CommentList []Comment
	var TempComment Comment
	token := c.Query("token")
	if token == "" {
		panic("没读到")
	}
	video_id := c.Query("video_id")
	if video_id == "" {
		panic("没读到")
	}
	rows, err := DB.Query("select id,content,createdate,comment_user_id from comment where video_id=" + video_id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&TempComment.Id, &TempComment.Content, &TempComment.CreateDate, &TempComment.User.Id)
		if err != nil {
			panic(err)
		}
		user_rows, err := DB.Query("select name,followcount,followercount from user where id=" + strconv.Itoa(int(TempComment.User.Id)))
		if err != nil {
			panic(err)
		}
		for user_rows.Next() {
			err = user_rows.Scan(&TempComment.User.Name, &TempComment.User.FollowCount, &TempComment.User.FollowerCount)
			if err != nil {
				panic(err)
			}
		}
		TempComment.User.IsFollow = false
		CommentList = append(CommentList, TempComment)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
		},
		CommentList: CommentList,
	})
}
