package handler

import (
	"beaconCloud/rpcServer/device-srv/db/model"

	"golang.org/x/net/context"

	"beaconCloud/rpcServer/device-srv/proto/device"
	"beaconCloud/rpcServer/device-srv/util"

	"github.com/micro/go-micro/errors"
	"github.com/pborman/uuid"
)

type Device struct{}

var beaconModel = new(model.Beacon)
var gateway = new(model.Gateway)
var beaconSetModel = new(model.BeaconSetting)

/* <<<<<<<<<<    Gateway Action Start ...   >>>>>>>>>>*/
/* 添加网关 */
func (s *Device) AddGateway(ctx context.Context, req *device.AddGatewayRequest, rsp *device.AddGatewayResponse) error {
	var e error
	rsp.Result = new(device.Result)
	rsp.Result.Id = "Handler AddGateway function"
	req.Gateway.ObjectId = uuid.NewUUID().String()
	if req.Gateway, e = util.VerifGateway(req.Gateway); e != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "网关基本数据校验失败!"
		rsp.Result.Status = "err"
		return e
	}

	if err := gateway.Add(req, rsp); err != nil {
		return err
	}
	return nil
}

/* 读取单个网关， 参数mac or object_id */
func (s *Device) ReadOneGateway(ctx context.Context, req *device.ReadOneGatewayRequest, rsp *device.ReadOneGatewayResponse) error {
	model.UseingDB(req.DatabaseUrl)
	if len(req.ObjectId) <= 0 && len(req.Mac) <= 0 {
		return errors.New("ReadOneGateway", "ReadOneGateway Function parameter is NULL", 400)
	}

	if len(req.ObjectId) > 0 && len(req.Mac) > 0 {
		err := gateway.ReadOneByObjectIdAndMac(req, rsp)
		return err
	}
	if len(req.ObjectId) > 0 {
		err := gateway.ReadOneByObjectId(req.ObjectId, rsp)
		return err
	}

	if len(req.Mac) > 0 {
		err := gateway.ReadOneByMac(req.Mac, rsp)
		return err
	}
	return nil
}

/* 读取全部网关 */
func (s *Device) ReadAllGateway(ctx context.Context, req *device.ReadAllGatewayRequest, rsp *device.ReadAllGatewayResponse) error {
	if err := gateway.ReadAll(req, rsp); err != nil {
		return err
	}
	return nil
}

/* 分页读取网关 */
func (s *Device) ReadPagingGateway(ctx context.Context, req *device.ReadPagingGatewayRequest, rsp *device.ReadPagingGatewayResponse) error {
	if err := gateway.ReadPaging(req, rsp); err != nil {
		return err
	}
	return nil
}

/* 删除网关 参数mac or object_id */
func (s *Device) DeleteGateway(ctx context.Context, req *device.DeleteGatewayRequest, rsp *device.DeleteGatewayResponse) error {
	rsp.Result = new(device.Result)
	rsp.Result.Id = "Handler DeleteGateway function"
	if err := gateway.Delete(req, rsp); err != nil {
		return err
	}
	return nil
}

/* 修改网关 */
func (s *Device) UpdateGateway(ctx context.Context, req *device.UpdateGatewayRequest, rsp *device.UpdateGatewayResponse) error {
	model.UseingDB(req.DatabaseUrl)
	var e error
	rsp.Result = new(device.Result)
	rsp.Result.Id = "Handler UpdateGateway function"

	if req.Gateway, e = util.VerifGateway(req.Gateway); e != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "网关基本数据校验失败!"
		rsp.Result.Status = "err"
		return e
	}

	if err := gateway.Update(req.Gateway, rsp); err != nil {
		return err
	}
	return nil
}

/* 更新网关状态 */
func (s *Device) UpdateGatewayStatus(ctx context.Context, req *device.UpdateGatewayStatusRequest, rsp *device.UpdateGatewayStatusResponse) error {
	rsp.Result = new(device.Result)
	rsp.Result.Id = "Handler UpdateGatewayStatus function"
	if len(req.ObjectId) > 0 {
		if err := gateway.UpdateStatus(req, rsp); err != nil {
			return err
		}
	}
	return nil
}

/* <<<<<<<<<<    Beacon Action Start ...   >>>>>>>>>>*/
func (s *Device) ReadOneBeacon(ctx context.Context, req *device.ReadOneBeaconRequest, rsp *device.ReadOneBeaconResponse) error {
	if len(req.ObjectId) <= 0 && len(req.Mac) <= 0 {
		return errors.New("ReadOneBeacon", "ReadOneBeacon Function parameter is NULL", 400)
	}

	if len(req.ObjectId) > 0 && len(req.Mac) > 0 {
		err := beaconModel.ReadByObjectIdAndMac(req, rsp)
		return err
	}

	if len(req.ObjectId) > 0 {
		err := beaconModel.ReadByObjectId(req, rsp)
		return err
	}

	if len(req.Mac) > 0 {
		err := beaconModel.ReadByMac(req, rsp)
		return err
	}

	return nil
}

func (s *Device) ReadAllBeacon(ctx context.Context, req *device.ReadAllBeaconRequest, rsp *device.ReadAllBeaconResponse) error {
	if err := beaconModel.ReadAll(req, rsp); err != nil {
		return err
	}
	return nil
}

func (s *Device) AddBeacon(ctx context.Context, req *device.AddBeaconRequest, rsp *device.AddBeaconResponse) error {
	var e error
	rsp.Result = &device.Result{Id: "Handler AddBeacon"}
	if req.Beacon, e = util.VerifBeacon(req.Beacon); e == nil {
		err := beaconModel.Add(req, rsp)
		return err
	}
	return e
}

func (s *Device) UpdateBeacon(ctx context.Context, req *device.UpdateBeaconRequest, rsp *device.UpdateBeaconResponse) error {
	var e error
	rsp.Result = new(device.Result)
	rsp.Result.Id = "Handler UpdateBeacon function"

	if req.Beacon, e = util.VerifBeacon(req.Beacon); e != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "Beacon基本数据校验失败!"
		rsp.Result.Status = "err"
		return e
	}

	if err := beaconModel.Update(req, rsp); err != nil {
		return err
	}
	return nil
}

/* 删除Beacon,可通过object_id或者mac作为删除条件 */
func (s *Device) DeleteBeacon(ctx context.Context, req *device.DeleteBeaconRequest, rsp *device.DeleteBeaconResponse) error {
	rsp.Result = &device.Result{Id: "Handler DeleteBeacon function."}
	if len(req.ObjectId) == 0 && len(req.Mac) == 0 {
		rsp.Result.Code = 404
		rsp.Result.Detail = "DeleteBeacon Function parameter is NULL"
		rsp.Result.Status = "err"
		return nil
	}
	if err := beaconModel.Delete(req, rsp); err != nil {
		return err
	}
	return nil
}

func (s *Device) ReadPagingBeacon(ctx context.Context, req *device.ReadPagingBeaconRequest, rsp *device.ReadPagingBeaconResponse) error {
	err := beaconModel.ReadPaging(req, rsp)
	return err
}

func (s *Device) UpdateBeaconApplyStatus(ctx context.Context, req *device.UpdateBeaconApplyStatusRequest, rsp *device.UpdateBeaconApplyStatusResponse) error {
	err := beaconModel.UpdateApplyStatus(req, rsp)
	return err
}

/* <<<<<<<<< Beacon Setting Start ... >>>>>>>>>>> */
func (s *Device) AddBeaconSetting(ctx context.Context, req *device.AddBeaconSettingRequest, rsp *device.AddBeaconSettingResponse) error {
	var e error
	rsp.Result = new(device.Result)
	rsp.Result.Id = "Handler AddBeaconSetting function"

	if req.BeaconSetting, e = util.VerifBeaconSetting(req.BeaconSetting); e != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "BeaconSetting 基本数据校验失败!"
		rsp.Result.Status = "err"
		return e
	}

	if err := beaconSetModel.Add(req, rsp); err != nil {
		return err
	}
	return nil
}

func (s *Device) ReadBeaconSetting(ctx context.Context, req *device.ReadBeaconSettingRequest, rsp *device.ReadBeaconSettingResponse) error {
	if len(req.ObjectId) <= 0 && len(req.Mac) <= 0 {
		return errors.New("ReadOneBeacon", "ReadBeaconSetting Function parameter is NULL", 400)
	}
	err := beaconSetModel.Read(req, rsp)
	return err
}

func (s *Device) ReadBeaconSetByObjectIdAndVersion(ctx context.Context, req *device.ReadBeaconSetByObjectIdAndVersionRequest, rsp *device.ReadBeaconSetByObjectIdAndVersionResponse) error {
	if len(req.ObjectId) <= 0 && req.Version < 0 {
		return errors.New("ReadOneBeacon", "ReadBeaconSetByObjectIdAndVersion Function parameter is NULL", 400)
	}
	err := beaconSetModel.ReadByObjectIdAndVersion(req, rsp)
	return err
}

func (s *Device) UpdateBeaconApplyVersion(ctx context.Context, req *device.UpdateBeaconApplyVersionRequest, rsp *device.UpdateBeaconApplyVersionResponse) error {
	if len(req.ObjectId) <= 0 && req.Version < 0 {
		return errors.New("ReadOneBeacon", "UpdateBeaconApplyVersion Function parameter is NULL", 400)
	}
	err := beaconSetModel.UpdateApplyVersion(req, rsp)
	return err
}

/* >>>>>>>>>>>>>>>>>>> InitDB >>>>>>>>>>>>>>>>>>>*/
func (s *Device) InitDB(ctx context.Context, req *device.InitDBRequest, rsp *device.InitDBResponse) error {
	rsp.Result = new(device.Result)
	if len(req.DatabaseUrl) > 0 {
		model.UseingDB(req.DatabaseUrl)
		rsp.Result.Id = "InitDB"
		rsp.Result.Code = 200
		rsp.Result.Detail = "InitDB OK"
		rsp.Result.Status = "ok"
		return nil
	}else {
		rsp.Result.Id = "InitDB"
		rsp.Result.Code = 500
		rsp.Result.Detail = "InitDB err"
		rsp.Result.Status = "err"
		return errors.New("InitDB", "DatabaseUrl is null", 404)
	}
}
