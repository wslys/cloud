package main

import (
	"beaconCloud/api/account"
	"beaconCloud/api/auth"
	"beaconCloud/api/device"
	"beaconCloud/rpcServer/user-srv/proto/loginCache"
	redis "beaconCloud/rpcServer/user-srv/redisCache"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-web"
	"golang.org/x/net/context"
	"log"
	"reflect"
	"strings"
	"time"
	"github.com/micro/go-micro/errors"
	"github.com/widuu/goini"
	"os"
)

var (
	basicDbUrl string
	conf *goini.Config
)

//检查元素在数组中是否存在
func Contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("Contains", "Not In", 404)
}

func judgeLogined(userId string) *loginCache.GetLoginResponse {
	login_cache := &loginCache.Login{
		UserId: userId,
	}

	gep := &loginCache.GetLoginRequest{
		Login: login_cache,
	}
	gsp := &loginCache.GetLoginResponse{
		Login: &loginCache.Login{},
	}

	redisCli := &redis.RedisCLI{}
	redisCli.GetLogin(context.TODO(), gep, gsp)

	if len(gsp.Login.UserId) > 0 {
		return gsp
	}

	return nil
}

// TODO  Global Filter
func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	var publicRoute = []string{"/account/createUser", "/account/login", "/auth/auth/{client_id}/{responseType}/{state}/{redirectUri}", "/auth/token"}
	userId := strings.TrimSpace(req.PathParameter("user_id"))
	routePath := req.SelectedRoutePath()

	ok, err := Contains(routePath, publicRoute)
	if err != nil {
		fmt.Println(err)
	}

	if ok == false {
		if user := judgeLogined(userId); user != nil {
			if user.Login.ValidityTime < time.Now().Unix() {
				login_cache := &loginCache.Login{
					UserId:userId,
				}
				dep := &loginCache.DeleteLoginRequest{
					Login:login_cache,
				}
				dsp := &loginCache.DeleteLoginResponse{}
				redisCli := &redis.RedisCLI{}
				redisCli.DeleteLogin(context.TODO(), dep, dsp);

				resp.WriteErrorString(402, "402: Login expired.登陆过期")
				return
			}
			dbUrl := basicDbUrl + user.Login.Username + "?charset=utf8"
			req.Request.Header.Set("db_url", dbUrl)
			chain.ProcessFilter(req, resp)
			return
		}
		// 登陆 / 授权验证
		resp.WriteError(401, errors.New("globalLogging", "globalLogging Login expired", 401))
		//resp.WriteErrorString(401, "401:")
		return
	} else {
		chain.ProcessFilter(req, resp)
	}
}

func init() {
	conf = goini.SetConfig("/etc/beacon-cloud/conf/conf.ini")
	username := conf.GetValue("database", "username")
	password := conf.GetValue("database", "password")
	hostname := conf.GetValue("database", "hostname")

	if len(username) == 0 || len(password) == 0 || len(hostname) == 0 {
		log.Fatal(errors.New("Device-SRV Init", "Config Database Username or Password or hostname IS NULL.", 500))
		os.Exit(1)
	}
	if len(username) == 0 || len(password) == 0 || len(hostname) == 0 {
		fmt.Println(errors.New("Api Init", "Config Database Username or Password or hostname IS NULL.", 500))
		os.Exit(1)
	}
	basicDbUrl = username + ":" + password + "@tcp("+ hostname +":3306)/beacon_cloud_"
}

func main() {
	service := web.NewService(
		web.Name("beaconcloud.api"),
		web.Address(":8081"),
	)
	service.Init()

	wc := restful.NewContainer()
	wc.Filter(globalLogging)

	// device
	device := new(device.Device)
	ws_d := new(restful.WebService)
	ws_d.Consumes(restful.MIME_JSON, restful.MIME_JSON)
	ws_d.Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws_d.Path("/device")
	// >>>>> beacon
	ws_d.Route(ws_d.POST("/beacon/add/{user_id}").To(device.AddBeacon))
	ws_d.Route(ws_d.GET("/beacon/readOne/{user_id}/{object_id}").To(device.ReadOneBeacon))
	ws_d.Route(ws_d.GET("/beacon/readAll/{user_id}").To(device.ReadAllBeacon))
	ws_d.Route(ws_d.GET("/beacon/readPaging/{user_id}/{size}/{currentPage}/{order}").To(device.ReadPagingBeacon))
	ws_d.Route(ws_d.GET("/beacon/delete/{user_id}/{object_id}/{mac}").To(device.DeleteBeacon))
	ws_d.Route(ws_d.POST("/beacon/update/{user_id}").To(device.UpdateBeacon))
	ws_d.Route(ws_d.GET("/beacon/updateApplyStatus/{user_id}/{object_id}").To(device.UpdateBeaconApplyStatus))

	// >>>>> beaconSetting
	ws_d.Route(ws_d.POST("/beaconSetting/add/{user_id}").To(device.AddBeaconSetting))
	ws_d.Route(ws_d.GET("/beaconSetting/read/{user_id}/{object_id}/{mac}").To(device.ReadBeaconSetting))
	ws_d.Route(ws_d.GET("/beaconSetting/readByObjectIdAndVersion/{user_id}/{object_id}/{version}").To(device.ReadBeaconSetByObjectIdAndVersion))
	ws_d.Route(ws_d.GET("/beaconSetting/updateApplyVersion/object_id/{user_id}/{object_id}/version/{version}").To(device.UpdateBeaconApplyVersion))

	// >>>>> gateway
	ws_d.Route(ws_d.GET("/gateway/readOne/{user_id}/{object_id}").To(device.ReadOneGateway))
	ws_d.Route(ws_d.GET("/gateway/readAll/{user_id}").To(device.ReadAllGateway))
	ws_d.Route(ws_d.GET("/gateway/readPaging/{user_id}/{mac}/{size}/{currentPage}/{order}").To(device.ReadPagingGateway))
	ws_d.Route(ws_d.POST("/gateway/add/{user_id}").To(device.AddGateway))
	ws_d.Route(ws_d.PUT("/gateway/update/{user_id}").To(device.UpdateGateway))
	ws_d.Route(ws_d.GET("/gateway/delete/{user_id}/{object_id}/{mac}").To(device.DeleteGateway))
	ws_d.Route(ws_d.GET("/gateway/updateStatus/{user_id}/{object_id}/{status}").To(device.UpdateGatewayStatus))
	wc.Add(ws_d)

	// account
	user := new(account.User)
	ws_u := new(restful.WebService)
	ws_u.Consumes(restful.MIME_JSON, restful.MIME_JSON)
	ws_u.Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws_u.Path("/account")

	// >>>>> account
	ws_u.Route(ws_u.POST("/createUser").To(user.CreateUser))
	ws_u.Route(ws_u.POST("/updateUser/{user_id}").To(user.UpdateUser))
	ws_u.Route(ws_u.GET("/readUserById/{user_id}/{Id}").To(user.ReadUserById))
	ws_u.Route(ws_u.POST("/login").To(user.Login))
	ws_u.Route(ws_u.GET("/logout/{user_id}/{sessionid}").To(user.Logout))
	ws_u.Route(ws_u.POST("/readAccountRoute/{user_id}").To(user.ReadAccountRoute))
	ws_u.Route(ws_u.POST("/updatePasswd/{user_id}").To(user.UpdatePasswd))
	ws_u.Route(ws_u.GET("/readSession/{user_id}/{sessionid}").To(user.ReadSession))
	wc.Add(ws_u)

	// auth
	auth := new(auth.User)
	ws_a := new(restful.WebService)
	ws_a.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws_a.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws_a.Path("/auth")

	// >>>>> auth
	ws_a.Route(ws_a.GET("/auth/{client_id}/{responseType}/{state}/{redirectUri}").To(auth.Authorize))
	ws_a.Route(ws_a.GET("/readAccount/{user_id}/{id}").To(auth.ReadAccount))
	ws_a.Route(ws_a.POST("/createAccount/{user_id}").To(auth.CreateAccount))
	ws_a.Route(ws_a.POST("/updateAccount/{user_id}").To(auth.UpdateAccount))
	ws_a.Route(ws_a.GET("/deleteAccount/{user_id}/{id}").To(auth.DeleteAccount))
	ws_a.Route(ws_a.POST("/searchAccount/{user_id}").To(auth.SearchAccount))
	ws_a.Route(ws_a.POST("/token").To(auth.Token))
	ws_a.Route(ws_a.POST("/revoke/{user_id}").To(auth.Revoke))
	ws_a.Route(ws_a.POST("/introspect/{user_id}").To(auth.Introspect))
	wc.Add(ws_a)

	// Register Handler
	service.Handle("/", wc)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
