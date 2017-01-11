package model

import (
	"beaconCloud/rpcServer/device-srv/proto/device"
	"fmt"
	"github.com/micro/go-micro/errors"
	"time"
)

type BeaconSetting struct {
	Id          int32
	ObjectId    string
	Mac         string
	Version     int32
	ApplyStatus int32
	Setting     string
	CreateAt    int64
	ApplyAt     int64
}

// 自定义表名
func (bs *BeaconSetting) TableName() string {
	return "beacon_setting"
}

// add
func (bs *BeaconSetting) Add(req *device.AddBeaconSettingRequest, rsp *device.AddBeaconSettingResponse) error {
	UseingDB(req.DatabaseUrl)
	beaconSet := &BeaconSetting{
		ObjectId:    req.BeaconSetting.ObjectId,
		Mac:         req.BeaconSetting.Mac,
		Version:     req.BeaconSetting.Version,
		ApplyStatus: req.BeaconSetting.ApplyStatus,
		Setting:     req.BeaconSetting.Setting,
		CreateAt:    req.BeaconSetting.CreateAt,
		ApplyAt:     req.BeaconSetting.ApplyAt,
	}

	rsp.Result.Id = "BeaconSetting model AddBeaconSetting func"
	lv, err := getLastVersion(beaconSet.ObjectId)
	beaconSet.Version = lv + 1
	_, err = o.Insert(beaconSet)
	if err != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "Get an BeaconSetting. mac:" + beaconSet.Mac
		rsp.Result.Status = "error"
		fmt.Println(err)
		return errors.New("model BeaconSetting add func", "create '"+beaconSet.ObjectId+"' fail. (创建'"+beaconSet.ObjectId+"'失败)", 201)
	}

	beacon := &Beacon{
		ObjectId: beaconSet.ObjectId,
	}
	if o.Read(beacon, "object_id") == nil {
		beacon.LastSettingVersion = beaconSet.Version
		beacon.UpdateAt = time.Now().Unix()

		if num, err := o.Update(beacon); err == nil {
			fmt.Println(num)
		}
	}

	rsp.Result.Code = 200
	rsp.Result.Detail = "New Insert an BeaconSetting. object_id:" + beaconSet.ObjectId
	rsp.Result.Status = "ok"
	return nil
}

func (bs *BeaconSetting) Read(req *device.ReadBeaconSettingRequest, rsp *device.ReadBeaconSettingResponse) error {
	UseingDB(req.DatabaseUrl)
	var beaconSets []*BeaconSetting
	total, err := o.QueryTable("beacon_setting").Filter("object_id", req.ObjectId).Filter("mac", req.Mac).All(&beaconSets)
	if err != nil {
		return err
	}
	for i := 0; i < len(beaconSets); i++ {
		gt := &device.BeaconSetting{}

		gt.Id = beaconSets[i].Id
		gt.ObjectId = beaconSets[i].ObjectId
		gt.Mac = beaconSets[i].Mac
		gt.Version = beaconSets[i].Version
		gt.ApplyStatus = beaconSets[i].ApplyStatus
		gt.Setting = beaconSets[i].Setting
		gt.CreateAt = beaconSets[i].CreateAt
		gt.ApplyAt = beaconSets[i].ApplyAt

		rsp.BeaconSettings = append(rsp.BeaconSettings, gt)
	}
	rsp.Total = total
	return nil
}

func (bs *BeaconSetting) ReadByObjectIdAndVersion(req *device.ReadBeaconSetByObjectIdAndVersionRequest, rsp *device.ReadBeaconSetByObjectIdAndVersionResponse) error {
	UseingDB(req.DatabaseUrl)
	beaconSet := &BeaconSetting{}
	rsp.BeaconSetting = &device.BeaconSetting{}

	err := o.QueryTable(beaconSet).Filter("object_id", req.ObjectId).Filter("version", req.Version).One(beaconSet)
	if err != nil {
		return err
	}

	rsp.BeaconSetting.Id = beaconSet.Id
	rsp.BeaconSetting.ObjectId = beaconSet.ObjectId
	rsp.BeaconSetting.Mac = beaconSet.Mac
	rsp.BeaconSetting.Version = beaconSet.Version
	rsp.BeaconSetting.ApplyStatus = beaconSet.ApplyStatus
	rsp.BeaconSetting.Setting = beaconSet.Setting
	rsp.BeaconSetting.CreateAt = beaconSet.CreateAt
	rsp.BeaconSetting.ApplyAt = beaconSet.ApplyAt
	return nil
}

func (bs *BeaconSetting) UpdateApplyVersion(req *device.UpdateBeaconApplyVersionRequest, rsp *device.UpdateBeaconApplyVersionResponse) error {
	UseingDB(req.DatabaseUrl)
	beaconSet := &BeaconSetting{
		ObjectId: req.ObjectId,
		Version:  req.Version,
	}

	if err := o.Read(beaconSet, "object_id", "version"); err != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "The '" + req.ObjectId + "' BeaconSetting Non-existent."
		rsp.Result.Status = "err"
		return err
	}

	beaconSet.ApplyStatus = 1
	if _, err := o.Update(beaconSet); err != nil {
		rsp.Result.Code = 202
		rsp.Result.Detail = "Update Update BeaconSetting Apply Version Status fail."
		rsp.Result.Status = "err"
		return err
	} else {
		rsp.Result.Code = 200
		rsp.Result.Detail = "Update Update BeaconSetting Apply Version Status Success."
		rsp.Result.Status = "ok"
		return nil
	}
}

/* get one beacon_setting last version */
func getLastVersion(objectId string) (int32, error) {
	beaconSet := &BeaconSetting{}

	err := o.QueryTable(beaconSet).Filter("object_id", objectId).OrderBy("-version").One(beaconSet) // 返回 QuerySeter
	if err != nil {
		return 0, err
	}
	return beaconSet.Version, nil
}
