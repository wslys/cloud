package account

import (
	users "beaconCloud/rpcServer/user-srv/proto/account"

	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
	"beaconCloud/rpcServer/device-srv/proto/device"
	"github.com/micro/go-micro/errors"
)

type User struct{}

type updatePassword struct {
	UserId          string
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
}

type loginData struct {
	Email    string
	Username string
	Password string
}

var (
	userCLI  users.AccountClient
	loginCLI users.LoginDbInfClient
)

func (u *User) ReadAccountRoute(req *restful.Request, rsp *restful.Response) {
	Id := req.PathParameter("Id")
	if Id == "" {
		rsp.WriteErrorString(402, "request args invalid")
	}
	readReq := &users.ReadAccountRequest{}
	readReq.Id = Id
	resp, err := loginCLI.ReadAccountRoute(context.TODO(), readReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	resp.NodeUrl = resp.NodeId + resp.Id
	rsp.WriteEntity(resp)
}

func (u *User) CreateUser(req *restful.Request, rsp *restful.Response) {
	usr := &users.User{}
	err := req.ReadEntity(usr)

	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	usr.NodeId = "0"
	request := &users.CreateRequest{
		User: usr,
	}
	resp, err := userCLI.Create(context.TODO(), request)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	device_srv_cl := device.NewDeviceClient("go.micro.srv.device", client.DefaultClient)
	dscReq := &device.InitDBRequest{
		DatabaseUrl:"root:root@tcp(127.0.0.1:3306)/beacon_cloud_" + usr.Username + "?charset=utf8",
	}
	dscRsp, err := device_srv_cl.InitDB(context.TODO(), dscReq)
	if err != nil {
		rspErr := errors.New(dscRsp.Result.Id, dscRsp.Result.Detail, dscRsp.Result.Code)
		rsp.WriteError(500, rspErr)
		return
	}
	// TODO
	rsp.WriteEntity(resp)
}

func (u *User) ReadUserById(req *restful.Request, rsp *restful.Response) {
	Id := req.PathParameter("Id")
	if Id == "" {
		rsp.WriteErrorString(402, "request args invalid")
	}
	request := &users.ReadRequest{}
	request.Id = Id
	resp, err := userCLI.Read(context.TODO(), request)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)
}

func (u *User) UpdateUser(req *restful.Request, rsp *restful.Response) {
	usr := &users.User{}
	err := req.ReadEntity(usr)

	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	request := &users.UpdateRequest{
		User: usr,
	}
	resp, err := userCLI.Update(context.TODO(), request)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(resp)
}

func (u *User) UpdatePasswd(req *restful.Request, rsp *restful.Response) {
	up := &updatePassword{}
	err := req.ReadEntity(up)
	fmt.Println(up)
	if err != nil {
		rsp.WriteErrorString(500, "request args invalid")
		return
	}
	up.UserId = strings.TrimSpace(up.UserId)
	up.OldPassword = strings.TrimSpace(up.OldPassword)
	up.NewPassword = strings.TrimSpace(up.NewPassword)
	up.ConfirmPassword = strings.TrimSpace(up.ConfirmPassword)

	if up.UserId == "" || up.NewPassword == "" || up.OldPassword == "" || up.ConfirmPassword == "" {
		rsp.WriteErrorString(500, "request args invalid")
		return
	}
	if up.NewPassword != up.ConfirmPassword {
		rsp.WriteErrorString(500, "两次输入的密码不一致")
		return
	}
	request := &users.UpdatePasswordRequest{}
	request.UserId = up.UserId
	request.OldPassword = up.OldPassword
	request.NewPassword = up.NewPassword
	request.ConfirmPassword = up.ConfirmPassword
	value, err := userCLI.UpdatePassword(context.TODO(), request)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(value)
}

func (u *User) Login(req *restful.Request, rsp *restful.Response) {
	ld := &loginData{}
	err := req.ReadEntity(ld)

	if err != nil {
		rsp.WriteError(500, err)
		return
	}

	ld.Email = strings.TrimSpace(ld.Email)
	ld.Username = strings.TrimSpace(ld.Username)
	ld.Password = strings.TrimSpace(ld.Password)
	if ld.Email == "" || ld.Password == "" {
		rsp.WriteErrorString(402, "request args invalid")
		return
	}

	request := &users.LoginRequest{
		Email:    ld.Email,
		Password: ld.Password,
	}
	value, err := userCLI.Login(context.TODO(), request)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}

	rsp.WriteEntity(value)
}

func (u *User) Logout(req *restful.Request, rsp *restful.Response) {
	sessionId := req.PathParameter("sessionid")
	if sessionId == "" {
		rsp.WriteErrorString(500, "request args invalid")
	}
	request := &users.LogoutRequest{}
	request.SessionId = sessionId
	value, err := userCLI.Logout(context.TODO(), request)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(value)
}

func (u *User) ReadSession(req *restful.Request, rsp *restful.Response) {
	sessionId := req.PathParameter("sessionid")
	if sessionId == "" {
		rsp.WriteErrorString(500, "request args invalid")
	}
	request := &users.ReadSessionRequest{}
	request.SessionId = sessionId
	value, err := userCLI.ReadSession(context.TODO(), request)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(value)
}

func init() {
	userCLI = users.NewAccountClient("go.micro.srv.user", client.DefaultClient)
	loginCLI = users.NewLoginDbInfClient("go.micro.srv.user", client.DefaultClient)
}
