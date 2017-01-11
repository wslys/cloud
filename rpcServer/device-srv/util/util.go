package util

import (
	"beaconCloud/rpcServer/device-srv/proto/device"
	"encoding/json"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
	"strings"
	"time"
)

// 获取API HTTP 请求的 POST数据,并转为相应的结构
func ReqPost(req *api.Request, b interface{}) error {
	data := req.Body

	d := []byte(data)
	err := json.Unmarshal(d, b)
	if err != nil {
		return errors.BadRequest("go.micro.api.gateway", "Api Post Data cannot be blank")
	}

	return nil
}

// 获取API HTTP 请求的 GET参数
func ReqGet(req *api.Request, key string) (string, error) {
	value, ok := req.Get[key]

	if !ok || len(value.Values) == 0 {
		return "", errors.BadRequest("go.micro.api.gateway", key+" cannot be blank")
	}

	return strings.Join(value.Values, " "), nil
}

// 验证Beacon数据结构是否正确
func VerifBeacon(beacon *device.Beacon) (*device.Beacon, error) {
	if beacon == nil {
		return nil, errors.New("Handler AddBeacon Verification", "AddBeacon Function parameter is NULL", 404)
	}

	if len(beacon.ObjectId) == 0 {
		return nil, errors.New("Handler AddBeacon Verification", "AddBeacon Function parameter object_id is Null", 404)
	}

	if len(beacon.UserId) == 0 {
		return nil, errors.New("Handler AddBeacon Verification", "AddBeacon Function parameter user_id is Null", 404)
	}

	if len(beacon.Mac) == 0 {
		return nil, errors.New("Handler AddBeacon Verification", "AddBeacon Function parameter mac is Null", 404)
	}

	beacon.ApplyStatus = 0

	if beacon.LastSettingVersion < 0 { // 最后一次编辑的配置项版本
		beacon.LastSettingVersion = 0
	}

	if beacon.CurrentVersion < 0 { // 当前启用的配置项版本
		beacon.CurrentVersion = 0
	}

	if beacon.Type <= 0 { // 1=ibeacon；
		beacon.Type = 1
	}

	if beacon.Status <= 0 { //  1=inactive(无效);2=active(启用);
		beacon.Status = 1
	}

	beacon.UpdateAt = time.Now().Unix()

	beacon.CreateAt = time.Now().Unix()

	return beacon, nil
}

// 验证Gateway数据结构是否正确
func VerifGateway(gateway *device.Gateway) (*device.Gateway, error) {
	if gateway == nil {
		return nil, errors.New("Util VerifGateway Verification", "AddGateway Function parameter is NULL", 404)
	}

	gateway.ObjectId = strings.TrimSpace(gateway.ObjectId)
	if len(gateway.ObjectId) == 0 {
		return nil, errors.New("Util VerifGateway Verification", "AddGateway Function parameter object_id is Null", 404)
	}

	gateway.Mac = strings.TrimSpace(gateway.Mac)
	if len(gateway.Mac) == 0 {
		return nil, errors.New("Util VerifGateway Verification", "AddGateway Function parameter mac is Null", 404)
	}

	gateway.Name = strings.TrimSpace(gateway.Name)
	if len(gateway.Name) == 0 {
		return nil, errors.New("Util VerifGateway Verification", "AddGateway Function parameter Name is Null", 404)
	}

	gateway.LastTime = time.Now().Unix()

	gateway.UpdateAt = time.Now().Unix()

	gateway.CreateAt = time.Now().Unix()

	return gateway, nil
}

// 验证BeaconSetting数据结构是否正确
func VerifBeaconSetting(beaconSet *device.BeaconSetting) (*device.BeaconSetting, error) {
	if beaconSet == nil {
		return nil, errors.New("Util VerifbeaconSet Verification", "AddbeaconSet Function parameter is NULL", 404)
	}

	beaconSet.ObjectId = strings.TrimSpace(beaconSet.ObjectId)
	if len(beaconSet.ObjectId) == 0 {
		return nil, errors.New("Util VerifbeaconSet Verification", "AddbeaconSet Function parameter object_id is Null", 404)
	}

	beaconSet.Mac = strings.TrimSpace(beaconSet.Mac)
	if len(beaconSet.Mac) == 0 {
		return nil, errors.New("Util VerifbeaconSet Verification", "AddbeaconSet Function parameter mac is Null", 404)
	}

	beaconSet.Setting = strings.TrimSpace(beaconSet.Setting)
	if len(beaconSet.Setting) == 0 {
		return nil, errors.New("Util VerifbeaconSet Verification", "AddbeaconSet Function parameter Setting is Null", 404)
	}

	beaconSet.CreateAt = time.Now().Unix()
	beaconSet.ApplyAt = time.Now().Unix()

	return beaconSet, nil
}
