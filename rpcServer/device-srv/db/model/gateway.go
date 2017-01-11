package model

import (
	"beaconCloud/rpcServer/device-srv/proto/device"
	"strings"

	"github.com/micro/go-micro/errors"
)

type Gateway struct {
	Id       int32
	ObjectId string
	Name     string
	Mac      string
	Status   int32
	LastTime int64
	Site     string
	CreateAt int64
	UpdateAt int64
}

// 自定义表名
func (g *Gateway) TableName() string {
	return "gateway"
}

// Add
func (g *Gateway) Add(req *device.AddGatewayRequest, rsp *device.AddGatewayResponse) error {
	UseingDB(req.DatabaseUrl)
	gw := &Gateway{
		ObjectId: req.Gateway.ObjectId,
		Name:     req.Gateway.Name,
		Mac:      req.Gateway.Mac,
		Status:   req.Gateway.Status,
		LastTime: req.Gateway.LastTime,
		Site:     req.Gateway.Site,
		CreateAt: req.Gateway.CreateAt,
		UpdateAt: req.Gateway.UpdateAt,
	}

	rsp.Result.Id = "model gateway add func"

	if created, _, err := o.ReadOrCreate(gw, "mac"); err == nil {
		if created {
			rsp.Result.Code = 200
			rsp.Result.Detail = "New Insert an Gateway. object_id:" + gw.ObjectId
			rsp.Result.Status = "ok"
			return nil
		} else {
			rsp.Result.Code = 201
			rsp.Result.Detail = "Get an Gateway. mac:" + gw.Mac
			rsp.Result.Status = "error"
			return errors.New("model gateway add func", "create '"+gw.ObjectId+"' fail. (创建'"+gw.ObjectId+"'失败)", 201)
		}
	}
	rsp.Result.Code = 500
	rsp.Result.Detail = "System error"
	rsp.Result.Status = "error"
	return errors.New("model gateway add func", "System error", 500)
}

// ReadOne
func (g *Gateway) ReadOneByMac(mac string, rsp *device.ReadOneGatewayResponse) error {
	gw := &Gateway{Mac: mac}
	rsp.Gateway = &device.Gateway{}

	if err := o.Read(gw, "mac"); err == nil {
		rsp.Gateway.Id = gw.Id
		rsp.Gateway.ObjectId = gw.ObjectId
		rsp.Gateway.Mac = gw.Mac
		rsp.Gateway.Name = gw.Name
		rsp.Gateway.Status = gw.Status
		rsp.Gateway.Site = gw.Site
		rsp.Gateway.LastTime = gw.LastTime
		rsp.Gateway.UpdateAt = gw.UpdateAt
		rsp.Gateway.CreateAt = gw.CreateAt
		return nil
	} else {
		return err
	}
}

// ReadOne
func (g *Gateway) ReadOneByObjectId(objectId string, rsp *device.ReadOneGatewayResponse) error {
	gw := &Gateway{ObjectId: objectId}

	rsp.Gateway = &device.Gateway{}

	if err := o.Read(gw, "object_id"); err == nil {
		rsp.Gateway.Id = gw.Id
		rsp.Gateway.ObjectId = gw.ObjectId
		rsp.Gateway.Mac = gw.Mac
		rsp.Gateway.Name = gw.Name
		rsp.Gateway.Status = gw.Status
		rsp.Gateway.Site = gw.Site
		rsp.Gateway.LastTime = gw.LastTime
		rsp.Gateway.UpdateAt = gw.UpdateAt
		rsp.Gateway.CreateAt = gw.CreateAt
		return nil
	} else {
		return err
	}
}

// ReadOne
func (g *Gateway) ReadOneByObjectIdAndMac(req *device.ReadOneGatewayRequest, rsp *device.ReadOneGatewayResponse) error {
	UseingDB(req.DatabaseUrl)
	gw := &Gateway{
		ObjectId: req.ObjectId,
		Mac:      req.Mac,
	}

	rsp.Gateway = &device.Gateway{}

	if err := o.Read(gw, "object_id", "mac"); err == nil {
		rsp.Gateway.Id = gw.Id
		rsp.Gateway.ObjectId = gw.ObjectId
		rsp.Gateway.Mac = gw.Mac
		rsp.Gateway.Name = gw.Name
		rsp.Gateway.Status = gw.Status
		rsp.Gateway.Site = gw.Site
		rsp.Gateway.LastTime = gw.LastTime
		rsp.Gateway.UpdateAt = gw.UpdateAt
		rsp.Gateway.CreateAt = gw.CreateAt
		return nil
	} else {
		return err
	}
}

// Update
func (g *Gateway) Update(gateway *device.Gateway, rsp *device.UpdateGatewayResponse) error {
	gw := &Gateway{ObjectId: gateway.ObjectId}

	if err := o.Read(gw, "object_id"); err != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "The '" + gateway.ObjectId + "' Gateway Non-existent."
		rsp.Result.Status = "err"
		return err
	}

	gw.Name = gateway.Name
	gw.Mac = gateway.Mac
	gw.UpdateAt = gateway.UpdateAt
	gw.Site = gateway.Site
	if _, err := o.Update(gw); err != nil {
		rsp.Result.Code = 202
		rsp.Result.Detail = "Update Gateway fail."
		rsp.Result.Status = "err"
		return err
	} else {
		rsp.Result.Code = 200
		rsp.Result.Detail = "Update Gateway Success."
		rsp.Result.Status = "ok"
		return nil
	}
}

// UpdateStatus
func (g *Gateway) UpdateStatus(req *device.UpdateGatewayStatusRequest, rsp *device.UpdateGatewayStatusResponse) error {
	UseingDB(req.DatabaseUrl)
	gw := &Gateway{ObjectId: req.ObjectId}

	if err := o.Read(gw, "object_id"); err != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "The '" + req.ObjectId + "' Gateway Non-existent."
		rsp.Result.Status = "err"
		return err
	}

	gw.Status = req.Status
	if _, err := o.Update(gw); err != nil {
		rsp.Result.Code = 202
		rsp.Result.Detail = "Update Gateway Status fail."
		rsp.Result.Status = "err"
		return err
	} else {
		rsp.Result.Code = 200
		rsp.Result.Detail = "Update Gateway Status Success."
		rsp.Result.Status = "ok"
		return nil
	}
}

// Delete
func (g *Gateway) Delete(req *device.DeleteGatewayRequest, rsp *device.DeleteGatewayResponse) error {
	UseingDB(req.DatabaseUrl)
	gw := &Gateway{
		ObjectId: req.ObjectId,
		Mac:      req.Mac,
	}

	if err := o.Read(gw, "object_id", "mac"); err != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "The '" + req.ObjectId + "' Gateway Non-existent."
		rsp.Result.Status = "err"
		return err
	}
	err := o.Begin()
	if err != nil {
		rsp.Result.Code = 202
		rsp.Result.Detail = "Delete Gateway fail."
		rsp.Result.Status = "err"
		return err
	}

	if num, err := o.Delete(gw, "object_id"); err != nil {
		rsp.Result.Code = 203
		rsp.Result.Detail = "Delete Gateway fail."
		rsp.Result.Status = "err"
		return err
	} else {
		if num != 1 {
			err = o.Rollback()
			rsp.Result.Code = 204
			rsp.Result.Detail = "Delete Gateway fail. Rollback"
			rsp.Result.Status = "err"
			return err
		} else {
			err = o.Commit()
			rsp.Result.Code = 200
			rsp.Result.Detail = "Delete Gateway Success"
			rsp.Result.Status = "ok"
			return err
		}
	}
}

// ReadPaging
func (g *Gateway) ReadPaging(req *device.ReadPagingGatewayRequest, rsp *device.ReadPagingGatewayResponse) error {
	UseingDB(req.DatabaseUrl)

	gateway := &Gateway{}
	qs := o.QueryTable(gateway)

	if len(req.Order) != 0 {
		h := strings.Split(req.Order, " ")

		if len(h) == 2 {
			if h[0] == "update_at" && h[1] == "desc" {
				qs = qs.OrderBy("-update_at")
			}
			if h[0] == "update_at" && h[1] == "asc" {
				qs = qs.OrderBy("update_at")
			}

			if h[0] == "create_at" && h[1] == "desc" {
				qs = qs.OrderBy("-create_at")
			}
			if h[0] == "create_at" && h[1] == "asc" {
				qs = qs.OrderBy("create_at")
			}
			if h[0] == "mac" && h[1] == "desc" {
				qs = qs.OrderBy("-mac")
			}
			if h[0] == "mac" && h[1] == "asc" {
				qs = qs.OrderBy("mac")
			}
			if h[0] == "object_id" && h[1] == "desc" {
				qs = qs.OrderBy("-object_id")
			}
			if h[0] == "object_id" && h[1] == "asc" {
				qs = qs.OrderBy("object_id")
			}
		}
	}

	if req.Limit > 0 && req.Offset >= 0 {
		qs = qs.Limit(req.Limit, (req.Offset-1)*req.Limit)
	} else {
		qs = qs.Limit(10, 0)
	}

	if len(req.Mac) > 0 {
		qs = qs.Filter("mac__iexact", req.Mac)
	}

	var gateways []*Gateway
	qs.All(&gateways)

	for i := 0; i < len(gateways); i++ {
		gt := &device.Gateway{}

		gt.Id = gateways[i].Id
		gt.ObjectId = gateways[i].ObjectId
		gt.Mac = gateways[i].Mac
		gt.Status = gateways[i].Status
		gt.Name = gateways[i].Name
		gt.CreateAt = gateways[i].CreateAt
		gt.LastTime = gateways[i].LastTime
		gt.Site = gateways[i].Site
		gt.UpdateAt = gateways[i].UpdateAt

		rsp.Gateways = append(rsp.Gateways, gt)
	}

	if cnt, ce := o.QueryTable(gateway).Count(); ce == nil {
		rsp.Total = cnt
	}

	return nil
}

// ReadAll
func (g *Gateway) ReadAll(req *device.ReadAllGatewayRequest, rsp *device.ReadAllGatewayResponse) error {
	UseingDB(req.DatabaseUrl)
	gateway := &Gateway{}
	var gateways []*Gateway
	_, err := o.QueryTable(gateway).All(&gateways)
	if err != nil {
		return err
	}
	for i := 0; i < len(gateways); i++ {
		gt := &device.Gateway{}

		gt.Id = gateways[i].Id
		gt.ObjectId = gateways[i].ObjectId
		gt.Mac = gateways[i].Mac
		gt.Status = gateways[i].Status
		gt.Name = gateways[i].Name
		gt.CreateAt = gateways[i].CreateAt
		gt.LastTime = gateways[i].LastTime
		gt.Site = gateways[i].Site
		gt.UpdateAt = gateways[i].UpdateAt

		rsp.Gateways = append(rsp.Gateways, gt)
	}
	if cnt, ce := o.QueryTable(gateway).Count(); ce != nil {
		rsp.Total = cnt
	}
	return nil
}
