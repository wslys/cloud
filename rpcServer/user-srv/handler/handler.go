package handler

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"crypto/sha512"

	"beaconCloud/rpcServer/user-srv/proto/loginCache"
	"github.com/micro/go-micro/errors"
	"beaconCloud/rpcServer/user-srv/db"
	account "beaconCloud/rpcServer/user-srv/proto/account"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"github.com/pborman/uuid"
	redis "beaconCloud/rpcServer/user-srv/redisCache"
	"fmt"
	"time"
)

const (
	HashMath       = 1
	HashVer        = 1
	HashAnswerFill = 0
	answerFill     = "#Answer#"
	x              = "cruft123"
)

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func random(i int) string {
	bytes := make([]byte, i)
	for {
		rand.Read(bytes)
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
		return string(bytes)
	}
	return "ughwhy?!!!"
}

func salt(ctx context.Context, answer string) (
	answerSalt string, err error) {
	r := make([]byte, 64+60)
	rand.Read(r)
	sha := sha512.New()
	sha.Write(r[20:84])
	sb := sha.Sum([]byte{HashMath, HashVer})
	answerSalt = base64.URLEncoding.EncodeToString(sb)
	return
}
func answerSalt(answer string) (answerSalt string) {
	for len(answerSalt) < 64 {
		if len(answerSalt) == 0 {
			answerSalt = answer + answerFill + answer
		} else {
			answerSalt = answerSalt + answerFill + answer
		}
	}
	return
}

type Account struct{}

func (s *Account) Create(ctx context.Context, req *account.CreateRequest, rsp *account.CreateResponse) error {
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.User.Password), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	req.User.Username = strings.ToLower(req.User.Username)
	req.User.Email = strings.ToLower(req.User.Email)
	if len(req.User.Id) != 0 {
		req.User.Id = ""
	}
	req.User.Id = uuid.NewUUID().String()
	return db.Create(req.User, salt, pp)

	if req.User.Username == "" || req.User.Password == "" {
		return errors.InternalServerError("go.micro.srv.user.Create", "account or passwd is null")
	}
	return nil
}

func (s *Account) Read(ctx context.Context, req *account.ReadRequest, rsp *account.ReadResponse) error {
	user, err := db.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

func (s *Account) Update(ctx context.Context, req *account.UpdateRequest, rsp *account.UpdateResponse) error {
	req.User.Username = strings.ToLower(req.User.Username)
	req.User.Email = strings.ToLower(req.User.Email)
	return db.Update(req.User)
}

func (s *Account) Delete(ctx context.Context, req *account.DeleteRequest, rsp *account.DeleteResponse) error {
	return db.Delete(req.Id)
}

func (s *Account) Search(ctx context.Context, req *account.SearchRequest, rsp *account.SearchResponse) error {
	users, err := db.Search(req.Username, req.Email, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}

func (s *Account) UpdatePassword(ctx context.Context, req *account.UpdatePasswordRequest, rsp *account.UpdatePasswordResponse) error {
	usr, err := db.Read(req.UserId)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	salt, hashed, err := db.SaltAndPassword(usr.Username, usr.Email)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.OldPassword)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.updatepassword", err.Error())
	}

	salt = random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	if err := db.UpdatePassword(req.UserId, salt, pp); err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}
	return nil
}

func (s *Account) Login(ctx context.Context, req *account.LoginRequest, rsp *account.LoginResponse) error {
	email := strings.ToLower(req.Email)

	salt, hashed, err := db.SaltAndPassword("", email)
	if err != nil {
		return err
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.Login", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.Password)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.login", err.Error())
	}
	u, err := db.ReadUserByUsernameAndEmail(email)
	if err!= nil {
		fmt.Println("DB ERROR :>>>>>>>>>>>. ", err)
	}

	// save session
	//sess := &account.Session{
	//	Id:       u.Id,
	//	Username: u.Username,
	//	Created:  time.Now().Unix(),
	//	Expires:  time.Now().Add(time.Hour * 24 * 7).Unix(),
	//}

	login_cache := &loginCache.Login{
		UserId:u.Id,
		NodeId:u.NodeId,
		Username:u.Username,
		Email:u.Email,
		Validity:6*60*60,
		ValidityTime:time.Now().Unix() + 6*60*60,
	}
	aeq := &loginCache.AddLoginRequest{
		Login:login_cache,
	}
	asp := &loginCache.AddLoginResponse{}

	gep := &loginCache.GetLoginRequest{
		Login:login_cache,
	}
	gsp := &loginCache.GetLoginResponse{
		Login:&loginCache.Login{},
	}

	redisCli := &redis.RedisCLI{}
	redisCli.GetLogin(context.TODO(), gep, gsp);

	if len(gsp.Login.UserId) > 0 {
		dep := &loginCache.DeleteLoginRequest{
			Login:login_cache,
		}
		dsp := &loginCache.DeleteLoginResponse{}
		redisCli.DeleteLogin(context.TODO(), dep, dsp);
	}
	redisCli.AddLogin(context.TODO(), aeq, asp);

	//if err := db.CreateSession(sess); err != nil {
	//	return errors.InternalServerError("go.micro.srv.user.Login", err.Error())
	//}
	rsp.User = &account.User{
		Id:u.Id,
		NodeId:u.NodeId,
		Username:u.Username,
		Email:u.Email,
		Created:u.Created,
		Updated:u.Updated,
	}
	return nil
}

func (s *Account) Logout(ctx context.Context, req *account.LogoutRequest, rsp *account.LogoutResponse) error {
	login_cache := &loginCache.Login{
		UserId:req.SessionId,
	}
	dep := &loginCache.DeleteLoginRequest{
		Login:login_cache,
	}
	dsp := &loginCache.DeleteLoginResponse{}
	redisCli := &redis.RedisCLI{}
	return redisCli.DeleteLogin(context.TODO(), dep, dsp);
	//return db.DeleteSession(req.SessionId)
}

func (s *Account) ReadSession(ctx context.Context, req *account.ReadSessionRequest, rsp *account.ReadSessionResponse) error {
	sess, err := db.ReadSession(req.SessionId)
	if err != nil {
		return err
	}
	rsp.Session = sess
	return nil
}

func (s *Account) RetrievePassword(ctx context.Context, req *account.RetrievePasswordRequest, rsp *account.RetrievePasswordResponse) error {
	return nil
}
