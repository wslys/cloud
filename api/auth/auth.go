package auth

import (
	account "beaconCloud/rpcServer/auth-srv/proto/account"
	oauth2 "beaconCloud/rpcServer/auth-srv/proto/oauth2"
	"strconv"
	"time"

	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type User struct{}

type caData struct {
	Id           string
	Type         string
	ClientId     string
	ClientSecret string
}

type tokenS struct {
	Code         string
	ClientId     string
	ClientSecret string
	RedirectUri  string
	GrantType    string
}

var (
	accountCLI account.AccountClient
	oauth2CLI  oauth2.Oauth2Client
)

func (user *User) ReadAccount(req *restful.Request, rsp *restful.Response) {
	Id := req.PathParameter("id")
	fmt.Println(Id)
	if Id == "" {
		rsp.WriteErrorString(402, "request args invalid")
	}
	readReq := &account.ReadRequest{
		Id: Id,
	}
	resp, err := accountCLI.Read(context.TODO(), readReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)
}

func (user *User) CreateAccount(req *restful.Request, rsp *restful.Response) {
	ca := &caData{}
	err := req.ReadEntity(ca)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}

	if ca.ClientId == "" || ca.ClientSecret == "" || ca.Type == "" {
		rsp.WriteErrorString(402, "request args invalid")
		return
	}
	createReq := &account.CreateRequest{
		Account: &account.Record{
			ClientId:     ca.ClientId,
			ClientSecret: ca.ClientSecret,
			Type:         ca.Type,
		},
	}
	//	createReq.Account = accounts
	createReq.Account.Created = time.Now().Unix()
	createReq.Account.Updated = time.Now().Unix()
	resp, err := accountCLI.Create(context.TODO(), createReq)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(resp)

}

func (user *User) UpdateAccount(req *restful.Request, rsp *restful.Response) {
	ca := &caData{}
	err := req.ReadEntity(ca)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}

	if ca.Id == "" || ca.ClientId == "" || ca.ClientSecret == "" || ca.Type == "" {
		rsp.WriteErrorString(402, "request args invalid")
		return
	}

	updateReq := &account.UpdateRequest{
		Account: &account.Record{
			Id:           ca.Id,
			ClientId:     ca.ClientId,
			ClientSecret: ca.ClientSecret,
			Type:         ca.Type,
		},
	}

	resp, err := accountCLI.Update(context.TODO(), updateReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)
}

func (user *User) DeleteAccount(req *restful.Request, rsp *restful.Response) {
	id := req.PathParameter("id")
	if id == "" {
		rsp.WriteErrorString(402, "request args invalid")
		return
	}
	deleteReq := &account.DeleteRequest{Id: id}
	resp, err := accountCLI.Delete(context.TODO(), deleteReq)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(resp)
}

func (user *User) SearchAccount(req *restful.Request, rsp *restful.Response) {
	client_id := req.PathParameter("client_id")
	accType := req.PathParameter("type")
	limit := req.PathParameter("limit")
	offset := req.PathParameter("offset")
	if client_id == "" || accType == "" || limit == "" || offset == "" {
		rsp.WriteErrorString(402, "request args invalid")
	}
	limits, err := strconv.Atoi(limit)
	if err != nil {
		rsp.WriteError(500, err)
	}
	offsets, err := strconv.Atoi(offset)
	if err != nil {
		rsp.WriteError(500, err)
	}
	searchReq := &account.SearchRequest{}
	searchReq.ClientId = client_id
	searchReq.Type = accType
	searchReq.Limit = int64(limits)
	searchReq.Offset = int64(offsets)
	resp, err := accountCLI.Search(context.TODO(), searchReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)
}

func (user *User) Authorize(req *restful.Request, rsp *restful.Response) {
	client_id := req.PathParameter("client_id")
	redirectUri := "https://wx02.yunli-wuli.com"
	responseType := req.PathParameter("responseType")
	state := req.PathParameter("state")
	//	scopes := req.PathParameter("scopes")
	if client_id == "" || redirectUri == "" || responseType == "" || state == "" {
		rsp.WriteErrorString(402, "request args invalid")
	}
	authReq := &oauth2.AuthorizeRequest{}
	authReq.ClientId = client_id
	authReq.RedirectUri = redirectUri
	authReq.ResponseType = responseType
	//	authReq.Scopes = scopes
	authReq.State = state

	resp, err := oauth2CLI.Authorize(context.TODO(), authReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)
}

func (user *User) Token(req *restful.Request, rsp *restful.Response) {
	ts := &tokenS{}
	err := req.ReadEntity(ts)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}

	if ts.ClientId == "" || ts.ClientSecret == "" || ts.GrantType == "" || ts.Code == "" || ts.RedirectUri == "" {
		fmt.Println(ts)
		rsp.WriteErrorString(402, "request args invalid")
		return
	}
	fmt.Println(ts)

	tokenReq := &oauth2.TokenRequest{
		ClientId:     ts.ClientId,
		ClientSecret: ts.ClientSecret,
		Code:         ts.Code,
		RedirectUri:  ts.RedirectUri,
		GrantType:    ts.GrantType,
	}

	resp, err := oauth2CLI.Token(context.TODO(), tokenReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)
}
func (user *User) Revoke(req *restful.Request, rsp *restful.Response) {
	access_token := req.PathParameter("access_token")
	refresh_token := req.PathParameter("refresh_token")
	if access_token == "" || refresh_token == "" {
		rsp.WriteErrorString(402, "request args invalid")
	}
	revokeReq := &oauth2.RevokeRequest{}
	revokeReq.AccessToken = access_token
	revokeReq.RefreshToken = refresh_token

	resp, err := oauth2CLI.Revoke(context.TODO(), revokeReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)
}
func (user *User) Introspect(req *restful.Request, rsp *restful.Response) {
	access_token := req.PathParameter("access_token")
	if access_token == "" {
		rsp.WriteErrorString(402, "request args invalid")
	}
	instrospectReq := &oauth2.IntrospectRequest{}
	instrospectReq.AccessToken = access_token

	resp, err := oauth2CLI.Introspect(context.TODO(), instrospectReq)
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(resp)

}

func init() {
	accountCLI = account.NewAccountClient("go.micro.srv.auth", client.DefaultClient)
	oauth2CLI = oauth2.NewOauth2Client("go.micro.srv.auth", client.DefaultClient)
}
