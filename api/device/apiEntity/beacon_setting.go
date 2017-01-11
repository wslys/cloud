package apiEntity

import (
	"beaconCloud/rpcServer/device-srv/proto/device"
	"strconv"
	"github.com/emicklei/go-restful"
	"golang.org/x/net/context"
	"strings"
)

type BeaconSetting struct{}

// 添加BeaconSetting (Beacon 配置)
func (s *BeaconSetting) AddBeaconSetting(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	beaconSetting := new(device.BeaconSetting)
	err := req.ReadEntity(beaconSetting)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	response, err := cl.AddBeaconSetting(context.TODO(), &device.AddBeaconSettingRequest{
		DatabaseUrl   : DatabaseUrl,
		BeaconSetting : beaconSetting,
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
	return
}

// 通过 (object_id,version)获取Beacon的配置信息
func (s *BeaconSetting) ReadBeaconSetting(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))
	mac := strings.TrimSpace(req.PathParameter("mac"))

	response, err := cl.ReadBeaconSetting(context.TODO(), &device.ReadBeaconSettingRequest{
		DatabaseUrl : DatabaseUrl,
		ObjectId    : objectId,
		Mac         : mac,
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
	return
}

// 通过 (object_id)获取Beacon的配置信息
func (s *BeaconSetting) ReadBeaconSetByObjectIdAndVersion(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))
	version := strings.TrimSpace(req.PathParameter("version"))
	vi, _ := strconv.ParseInt(version, 10, 64)

	response, err := cl.ReadBeaconSetByObjectIdAndVersion(context.TODO(), &device.ReadBeaconSetByObjectIdAndVersionRequest{
		DatabaseUrl : DatabaseUrl,
		ObjectId    : objectId,
		Version     : int32(vi),
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
	return
}

// 修改 Beacon 应用的version
func (s *BeaconSetting) UpdateBeaconApplyVersion(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))
	version := strings.TrimSpace(req.PathParameter("version"))
	vi, _ := strconv.ParseInt(version, 10, 64)

	response, err := cl.UpdateBeaconApplyVersion(context.TODO(), &device.UpdateBeaconApplyVersionRequest{
		DatabaseUrl : DatabaseUrl,
		ObjectId    : objectId,
		Version     : int32(vi),
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
	return
}
