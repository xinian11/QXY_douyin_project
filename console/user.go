package console

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}
type CurrentUser struct {
	User     User
	password string
}

var Current_user CurrentUser

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	exist, err := DB.Query("select name,password from user where name=\"" + username + "\" and password=\"" + password + "\"")
	if err != nil {
		if err != nil {
			panic(err)
		}
	}
	if exist.Next() {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		user_insect, err := DB.Begin()
		if err != nil {
			panic(err)
		}
		r, err := DB.Exec("INSERT INTO user(name,password,token)VALUES(?,?,?)", username, password, token)
		if err != nil {
			panic(err)
		}
		user_id, err := r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			user_insect.Rollback()
			return
		}
		user_insect.Commit()
		Current_user.password = password
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user_id,
			Token:    username + password,
		})
	}

}
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	var Temp_user User
	Temp_user.Id = 0
	exist, err := DB.Query("select id,name,followcount,followercount from user where name=\"" + username + "\" and password=\"" + password + "\"")
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
		Current_user.password = password
	}
	if Temp_user.Id != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   Temp_user.Id,
			Token:    username + password,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	// token := c.Query("token")
	var Temp_user User
	Temp_user.Id = 0
	id := c.Query("user_id")
	if id == "" {
		panic("没读到")
	}
	exist, err := DB.Query("select id,name,followcount,followercount from user where id=" + id)
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
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     Temp_user,
		})
	}

}
