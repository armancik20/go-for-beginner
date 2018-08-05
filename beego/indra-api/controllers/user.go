package controllers

import (
	"github.com/indraoct/go-for-beginner/beego/indra-api/models"
	"encoding/json"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}


// @Title CreateUser
// @Description create users
// @Param	username		form 	string	true		"The username for signup"
// @Param	password		form 	string	true		"The password for signup"
// @Success 200 {object} models.ResponseInsertUser
// @router /add [post]
func (this *UserController) Add() {

	var users models.Users
	var resInsert models.ResponseInsertUser
	resInsert.Status = 0
	resInsert.Message = "failed to insert data user"

	username := this.Ctx.Request.PostForm.Get("username")
	password := this.Ctx.Request.PostForm.Get("password")

	users.Username = username
	users.Password = password

	id_user := models.AddUser(users)

	if(id_user != 0){
		resInsert.Status = 1
		resInsert.Message = "Success"
		resInsert.Data.Id_user = id_user
	}

	this.Data["json"] = resInsert
	this.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param	filter[username]		query 	string	false	"username"
// @Param	offset			query 	int	false	offset""
// @Param	limit			query 	int	false	"limit"
// @Param	filter[status]			query 	int	false	"status"
// @Success 200 {object} models.ResponseGetAllUser true  "body for user content"
// @router /getall [get]
func (this *UserController) GetAll() {
	var filter models.UserFilter
	var resUserGetAll models.ResponseGetAllUser

	resUserGetAll.Status = 0
	resUserGetAll.Message = "Data not found"
	if(this.Ctx.Input.IsGet() == true) {

		offset := this.GetString("offset")
		limit := this.GetString("limit")
		filter_status := this.GetString("filter[status]")
		filter_username := this.GetString("filter[username]")

		if ( filter_status != "") {
			filter.Status = filter_status
		}else{
			filter.Status = ""
		}

		if (filter_username != "") {
			filter.Username = filter_username
		}

		resUser := models.GetAllUsers(offset, limit, filter)

		if (resUser.Count > 0) {
			resUserGetAll.Status = 1
			resUserGetAll.Message = "Success"
			resUserGetAll.Data = resUser.Data
			resUserGetAll.Limit = resUser.Limit
			resUserGetAll.Offset = resUser.Offset
			resUserGetAll.Count = resUser.Count

		}

	}else{
		this.Ctx.Output.Status = 401
		resUserGetAll.Status = 0
		resUserGetAll.Message = "Invalid Method"
	}
	this.Data["json"] = resUserGetAll
	this.ServeJSON()
}

// @Title Get
// @Description get user by username
// @Param	username		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ResponseGetUser
// @Failure 403 :username is empty
// @router /get/:username [get]
func (this *UserController) Get() {
	var resUserGetUser models.ResponseGetUser
	resUserGetUser.Status = 0
	resUserGetUser.Message = "Data not found"
	if(this.Ctx.Input.IsGet() == true) {
		username := this.GetString(":username")
		if(username != "") {
			resUser := models.GetUser(username)
			if (resUser.Count > 0) {
				resUserGetUser.Status = 1
				resUserGetUser.Message = "Success"
				resUserGetUser.Data = resUser.Data
				resUserGetUser.Limit = resUser.Limit
				resUserGetUser.Offset = resUser.Offset
				resUserGetUser.Count = resUser.Count
			}
		}

	}else{
		this.Ctx.Output.Status = 401
		resUserGetUser.Status = 0
		resUserGetUser.Message = "Invalid Method"
	}
	this.Data["json"] = resUserGetUser
	this.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var users models.Users
		json.Unmarshal(u.Ctx.Input.RequestBody, &users)
		uu, err := models.UpdateUser(uid, &users)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		form 	string	true		"The username for login"
// @Param	password		form 	string	true		"The password for login"
// @Success 200 { "message": "login success", "status": 1, "tokensting": "string" }
// @Failure 403 { "message": "user not exist", "status": 0, "tokensting": "string" }
// @Failure 401 Invalid method
// @router /login [post]
func (this *UserController) Login() {
	if(this.Ctx.Input.IsPost() == true) {
		var resLogin models.ResponseLogin
		resLogin.Status = 0
		username := this.Ctx.Request.PostForm.Get("username")
		password := this.Ctx.Request.PostForm.Get("password")

		if models.Login(username, password) {
			resLogin.Status = 1
			resLogin.Message = "login success"
		} else {
			this.Ctx.Output.Status = 403
			resLogin.Message = "user not exist"
		}
		this.Data["json"] = resLogin
	}else{
		this.Ctx.Output.Status = 401
		this.Data["json"] = "invalid method!"

	}

	this.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}