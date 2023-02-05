package console

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	var Temp_user User
	token := c.Query("token")
	if token == "" {
		panic("没读到")
	}
	to_user_id := c.Query("to_user_id")
	if to_user_id == "" {
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
	follow_change, err := DB.Begin()
	if err != nil {
		panic(err)
	}
	if action_type == "1" {
		r, err := DB.Exec("update user set followercount=followercount+1 where id=" + to_user_id)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			follow_change.Rollback()
			return
		}
		r, err = DB.Exec("insert into follow(check_id,follow_user_id)values(?,?)", Temp_user.Id, to_user_id)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			follow_change.Rollback()
			return
		}
		follow_change.Commit()
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "关注成功",
		})
	}
	if action_type == "2" {
		r, err := DB.Exec("update user set followercount=followercount-1 where id=" + to_user_id)
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			follow_change.Rollback()
			return
		}
		r, err = DB.Exec("delete from follow where follow_user_id=" + to_user_id + " and check_id=" + strconv.Itoa(int(Temp_user.Id)))
		if err != nil {
			panic(err)
		}
		_, err = r.LastInsertId()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			follow_change.Rollback()
			return
		}
		follow_change.Commit()
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取消关注",
		})
	}
}

func RelationFollowList(c *gin.Context) {
	var Current_User_Follow []User
	var Temp_User User
	var Follow_User_Id []int
	var Temp_id int
	user_id := c.Query("user_id")
	rows, err := DB.Query("select follow_user_id from follow where check_id=" + user_id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&Temp_id)
		if err != nil {
			panic(err)
		}
		Follow_User_Id = append(Follow_User_Id, Temp_id)
	}
	for _, i := range Follow_User_Id {
		rows, err = DB.Query("select name,followcount,followercount from user where id=" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			err = rows.Scan(&Temp_User.Name, &Temp_User.FollowCount, &Temp_User.IsFollow)
			if err != nil {
				panic(err)
			}
			Temp_User.Id = int64(i)
			Temp_User.Name = Current_user.User.Name
			Temp_User.IsFollow = true
			Current_User_Follow = append(Current_User_Follow, Temp_User)
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: Current_User_Follow,
	})

}
func RelationFollowerList(c *gin.Context) {
	var Current_User_Follower []User
	var Temp_User User
	var Follower_User_Id []int
	var Temp_id int
	user_id := c.Query("user_id")
	rows, err := DB.Query("select check_id from follow where follow_user_id=" + user_id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&Temp_id)
		if err != nil {
			panic(err)
		}
		Follower_User_Id = append(Follower_User_Id, Temp_id)
	}
	for _, i := range Follower_User_Id {
		rows, err = DB.Query("select name,followcount,followercount from user where id=" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			err = rows.Scan(&Temp_User.Name, &Temp_User.FollowCount)
			if err != nil {
				panic(err)
			}
			Temp_User.Id = int64(i)
			Temp_User.Name = Current_user.User.Name
			Current_User_Follower = append(Current_User_Follower, Temp_User)
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: Current_User_Follower,
	})

}

func RelationFriendList(c *gin.Context) {
	var Current_User_Follower []User
	var Temp_User User
	var Follower_User_Id []int
	var Temp_id int
	user_id := c.Query("user_id")
	rows, err := DB.Query("select check_id from follow where follow_user_id=" + user_id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&Temp_id)
		if err != nil {
			panic(err)
		}
		Follower_User_Id = append(Follower_User_Id, Temp_id)
	}
	for _, i := range Follower_User_Id {
		rows, err = DB.Query("select name,followcount,followercount from user where id=" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			err = rows.Scan(&Temp_User.Name, &Temp_User.FollowCount)
			if err != nil {
				panic(err)
			}
			Temp_User.Id = int64(i)
			Temp_User.Name = Current_user.User.Name
			Current_User_Follower = append(Current_User_Follower, Temp_User)
		}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: Current_User_Follower,
	})

}
