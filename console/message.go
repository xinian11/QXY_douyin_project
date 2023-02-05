package console

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

func MessageAction(c *gin.Context) {
	var Cur_user User
	var To_user User
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	Cur_user.Id = 0
	To_user.Id = 0
	// 判断当前用户是否登录
	exist, err := DB.Query("select id,name,followcount,followercount from user where token=\"" + token + "\"")
	if err != nil {
		panic(err)
	}
	for exist.Next() {
		err = exist.Scan(&Cur_user.Id, &Cur_user.Name, &Cur_user.FollowCount, &Cur_user.FollowerCount)
		if err != nil {
			panic(err)
		}
		Cur_user.IsFollow = false
		Current_user.User = Cur_user
	}
	if Cur_user.Id == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	//判断目标用户是否存在
	exist, err = DB.Query("select id,name,followcount,followercount from user where id=" + toUserId)
	if err != nil {
		panic(err)
	}
	for exist.Next() {
		err = exist.Scan(&To_user.Id, &To_user.Name, &To_user.FollowCount, &To_user.FollowerCount)
		if err != nil {
			panic(err)
		}
		To_user.IsFollow = false
		Current_user.User = Cur_user
	}
	if To_user.Id != 0 {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(Cur_user.Id, int64(userIdB))
		if len(content) > 100 {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "聊天字数过长。"})
			return
		}
		message_save, err := DB.Begin()
		if err != nil {
			panic(err)
		}
		createTime := time.Now().Format(time.Kitchen)
		r, err := DB.Exec("insert into message(chatkey,message,createtime)values(?,?,?)", chatKey, content, createTime)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			message_save.Rollback()
			return
		}
		message_save.Commit()
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	var Cur_user User
	var Current_message Message
	var MessageList []Message
	exist, err := DB.Query("select id,name,followcount,followercount from user where token=\"" + token + "\"")
	if err != nil {
		panic(err)
	}
	for exist.Next() {
		err = exist.Scan(&Cur_user.Id, &Cur_user.Name, &Cur_user.FollowCount, &Cur_user.FollowerCount)
		if err != nil {
			panic(err)
		}
		Cur_user.IsFollow = false
		Current_user.User = Cur_user

	}
	if Cur_user.Id != 0 {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(Cur_user.Id, int64(userIdB))
		rows, err := DB.Query("select id,message,createtime from message where chatkey=\"" + chatKey + "\"")
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			err = rows.Scan(&Current_message.Id, &Current_message.Content, &Current_message.CreateTime)
			if err != nil {
				panic(err)
			}
			MessageList = append(MessageList, Current_message)
		}
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: MessageList})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
