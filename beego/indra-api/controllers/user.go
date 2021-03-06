package controllers

import (
	"github.com/indraoct/go-for-beginner/beego/indra-api/models"
	"encoding/json"
	"strconv"
	"math"
)

// Operations about Users
type UserController struct {
	BaseController
}

//struct for Response Get User
type ResponseGetUser struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Count int `json:"count"`
	Data  []models.Users `json:"data"`
}

//struct for Response Get All User
type ResponseGetAllUser struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Page int `json:"page"`
	Perpage int `json:"perpage"`
	Count int `json:"count"`
	Data  []models.Users `json:"data"`
	Pages PagesList `json:"pages"`
}

type PagesList struct{
	First string `json:"first"`
	Last string `json:"last"`
}

type ResponseInsertUser struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data  struct{Id_user int64 `json:"id_user"`} `json:"data"`
}

//struct for Response Login
type ResponseLogin struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Tokensting string `json:"tokensting"`
}


// @Title CreateUser
// @Description create users
// @Param	username		form 	string	true		"The username for signup"
// @Param	password		form 	string	true		"The password for signup"
// @Success 200 { "status": 1, "message": "Success", "data": { "id_user": 4 } }
// @router /add [post]
func (this *UserController) Add() {

	var users models.Users
	var resInsert ResponseInsertUser
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
// @Param	page			query 	int	false	page""
// @Param	perpage			query 	int	false	"perpage"
// @Param	filter[status]			query 	int	false	"status"
// @Success 200 {"status":1,"message":"Success","page":1,"perpage":2,"count":3,"data":[{"id":2,"username":"dewibeego","password":"","status":1,"created_date":"2018-08-05 18:41:27","updated_date":"2018-08-05 18:41:27"},{"id":7,"username":"dewibeego2","password":"","status":1,"created_date":"2018-08-06 17:58:16","updated_date":"2018-08-06 17:58:16"}],"pages":{"first":"http://indraapi.dev/v1/user/getall/?filter[status]=1&filter[username]=dewi&page=1&perpage=2","last":"http://indraapi.dev/v1/user/getall/?filter[status]=1&filter[username]=dewi&page=2&perpage=2"}} true  "body for user content"
// @router /getall [get]
func (this *UserController) GetAll() {
	var filter models.UserFilter
	var resUserGetAll ResponseGetAllUser

	resUserGetAll.Status = 0
	resUserGetAll.Message = "Data not found"
	if(this.Ctx.Input.IsGet() == true) {

		page := this.GetString("page")
		perpage := this.GetString("perpage")
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

		resUser := models.GetAllUsers(page, perpage, filter)

		if (resUser.Count > 0) {
			resUserGetAll.Status = 1
			resUserGetAll.Message = "Success"
			resUserGetAll.Data = resUser.Data
			resUserGetAll.Perpage = resUser.Perpage
			resUserGetAll.Page = resUser.Page
			resUserGetAll.Count = resUser.Count

			//List Pages LOGIC
			prefix_url := "http://"+this.Ctx.Input.Host()+""+this.Ctx.Input.URL()
			first_page := prefix_url+"/?"+"filter[status]="+filter_status+"&filter[username]="+filter_username+"&page=1&perpage="+strconv.Itoa(resUser.Perpage)
			page_end := int(math.Round(float64(resUser.Count)/float64(resUser.Perpage)))
			last_page := prefix_url+"/?filter[status]="+filter_status+"&filter[username]="+filter_username+"&page="+strconv.Itoa(page_end)+"&perpage="+strconv.Itoa(resUser.Perpage)

			resUserGetAll.Pages.First = first_page
			resUserGetAll.Pages.Last = last_page

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
// @Success 200 { "status": 1, "message": "Success", "offset": 0, "limit": 0, "count": 1, "data": [ { "created_date": "2018-08-02 12:27:04", "status": "1", "updated_date": null, "username": "indrabeego" } ] }
// @Failure 403 :username is empty
// @router /get/:username [get]
func (this *UserController) Get() {
	var resUserGetUser ResponseGetUser
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
// @Param	username	form 	string	true		"The username for login"
// @Param	password	form 	string	true		"The password for login"
// @Success 200 { "message": "login success", "status": 1, "tokensting": "string" }
// @Failure 403 { "message": "user not exist", "status": 0, "tokensting": "string" }
// @Failure 401 Invalid method
// @router /login [post]
func (this *UserController) Login() {
		var resLogin ResponseLogin
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