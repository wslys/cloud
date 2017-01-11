package apiEntity

import (
	"strconv"

	"beaconCloud/rpcServer/device-srv/proto/device"
	"github.com/emicklei/go-restful"
	"golang.org/x/net/context"
	"strings"
)

type Gateway struct{}

// 读取单条网关信息
func (s *Gateway) ReadOneGateway(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))

	response, err := cl.ReadOneGateway(context.TODO(), &device.ReadOneGatewayRequest{
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

// 读取全部网关信息
func (s *Gateway) ReadAllGateway(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	response, err := cl.ReadAllGateway(context.TODO(), &device.ReadAllGatewayRequest{
		DatabaseUrl: DatabaseUrl,
	})
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(response)
}

// 添加网关信息
func (s *Gateway) AddGateway(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	gateway := new(device.Gateway)
	err := req.ReadEntity(gateway)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	response, err := cl.AddGateway(context.TODO(), &device.AddGatewayRequest{
		DatabaseUrl : DatabaseUrl,
		Gateway     : gateway,
	})
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(response)
}

// 分页读取网关信息
func (s *Gateway) ReadPagingGateway(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	size := strings.TrimSpace(req.PathParameter("size"))
	currentPage := strings.TrimSpace(req.PathParameter("currentPage"))
	order := strings.TrimSpace(req.PathParameter("order"))
	mac := strings.TrimSpace(req.PathParameter("mac"))

	limit, _ := strconv.ParseInt(size, 10, 64)
	offset, _ := strconv.ParseInt(currentPage, 10, 64)

	response, err := cl.ReadPagingGateway(context.TODO(), &device.ReadPagingGatewayRequest{
		DatabaseUrl: DatabaseUrl,
		Limit  : int32(limit),
		Offset : int32(offset),
		Order  : order,
		Mac    : mac,
	})
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(response)
}

// 修改网关信息
func (s *Gateway) UpdateGateway(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	gateway := new(device.Gateway)
	err := req.ReadEntity(gateway)
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	response, err := cl.UpdateGateway(context.TODO(), &device.UpdateGatewayRequest{
		DatabaseUrl : DatabaseUrl,
		Gateway     : gateway,
	})
	if err != nil {
		rsp.WriteError(500, err)
	}
	rsp.WriteEntity(response)
}

// 删除网关信息
func (s *Gateway) DeleteGateway(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))
	mac := strings.TrimSpace(req.PathParameter("mac"))

	response, err := cl.DeleteGateway(context.TODO(), &device.DeleteGatewayRequest{
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

// 修改网关状态
func (s *Gateway) UpdateGatewayStatus(req *restful.Request, rsp *restful.Response) {
	DatabaseUrl = req.Request.Header.Get("db_url")
	objectId := strings.TrimSpace(req.PathParameter("object_id"))
	status := strings.TrimSpace(req.PathParameter("status"))
	st, _ := strconv.ParseInt(status, 10, 64)
	response, err := cl.UpdateGatewayStatus(context.TODO(), &device.UpdateGatewayStatusRequest{
		DatabaseUrl : DatabaseUrl,
		ObjectId    : objectId,
		Status      : int32(st),
	})
	if err != nil {
		rsp.WriteError(500, err)
		return
	}
	rsp.WriteEntity(response)
	return
}
