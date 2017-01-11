package apiEntity

import (
	"beaconCloud/rpcServer/device-srv/proto/device"
	"github.com/emicklei/go-restful"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

type Beacon struct{}

// 添加Beacon
func (s *Beacon) AddBeacon(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	beacon := new(device.Beacon)
	err := req.ReadEntity(beacon)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	response, err := cl.AddBeacon(context.TODO(), &device.AddBeaconRequest{
		DatabaseUrl : DatabaseUrl,
		Beacon      : beacon,
	})
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(response)
}

// 读取一个beacon (根据object_id)
func (s *Beacon) ReadOneBeacon(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))
	//mac := strings.TrimSpace(req.PathParameter("mac"))

	response, err := cl.ReadOneBeacon(context.TODO(), &device.ReadOneBeaconRequest{
		DatabaseUrl : DatabaseUrl,
		ObjectId    : objectId,
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
}

// 读取所有的Beacon
func (s *Beacon) ReadAllBeacon(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")

	response, err := cl.ReadAllBeacon(context.TODO(), &device.ReadAllBeaconRequest{
		DatabaseUrl: DatabaseUrl,
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
	return
}

// 分页读取Beacon
func (s *Beacon) ReadPagingBeacon(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	size := strings.TrimSpace(req.PathParameter("size"))
	currentPage := strings.TrimSpace(req.PathParameter("currentPage"))
	order := strings.TrimSpace(req.PathParameter("order"))

	limit, _ := strconv.ParseInt(size, 10, 64)
	offset, _ := strconv.ParseInt(currentPage, 10, 64)

	response, err := cl.ReadPagingBeacon(context.TODO(), &device.ReadPagingBeaconRequest{
		DatabaseUrl: DatabaseUrl,
		Limit  : int32(limit),
		Offset : int32(offset),
		Order  : order,
	})
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(response)
}

// 删除Beacon (object_id)
func (s *Beacon) DeleteBeacon(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))
	mac := strings.TrimSpace(req.PathParameter("mac"))

	response, err := cl.DeleteBeacon(context.TODO(), &device.DeleteBeaconRequest{
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

// 更新Beacon
func (s *Beacon) UpdateBeacon(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	beacon := new(device.Beacon)
	err := req.ReadEntity(beacon)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	response, err := cl.UpdateBeacon(context.TODO(), &device.UpdateBeaconRequest{
		DatabaseUrl : DatabaseUrl,
		Beacon      : beacon,
	})
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(response)
}

// 更新Beacon 是否已同步到设备
func (s *Beacon) UpdateBeaconApplyStatus(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))

	response, err := cl.UpdateBeaconApplyStatus(context.TODO(), &device.UpdateBeaconApplyStatusRequest{
		DatabaseUrl : DatabaseUrl,
		ObjectId    : objectId,
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
	return
}
