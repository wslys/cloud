package model

import (
	"beaconCloud/rpcServer/device-srv/proto/device"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/micro/go-micro/errors"
)

type Beacon struct {
	Id                 int32
	ObjectId           string
	UserId             string
	Mac                string
	Status             int32
	Password           string
	ChangePassword     string
	Type               int32
	CurrentVersion     int32
	LastSettingVersion int32
	ApplyStatus        int32
	CreateAt           int64
	UpdateAt           int64
}

/* 读取所有Beacon */
func (b *Beacon) ReadAll(req *device.ReadAllBeaconRequest, rsp *device.ReadAllBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	var beacons []*Beacon
	total, err := o.QueryTable("beacon").All(&beacons)

	if err != nil {
		return err
	}

	for i := 0; i < len(beacons); i++ {
		bc := &device.Beacon{}

		bc.Id = beacons[i].Id
		bc.ObjectId = beacons[i].ObjectId
		bc.UserId = beacons[i].UserId
		bc.Mac = beacons[i].Mac
		bc.Status = beacons[i].Status
		bc.Password = beacons[i].Password
		bc.ChangePassword = beacons[i].ChangePassword
		bc.Type = beacons[i].Type
		bc.CurrentVersion = beacons[i].CurrentVersion
		bc.LastSettingVersion = beacons[i].LastSettingVersion
		bc.ApplyStatus = beacons[i].ApplyStatus
		bc.CreateAt = beacons[i].CreateAt
		bc.UpdateAt = beacons[i].UpdateAt

		rsp.Beacons = append(rsp.Beacons, bc)
	}
	rsp.Total = total
	return nil
}

func (b *Beacon) ReadByObjectId(req *device.ReadOneBeaconRequest, rsp *device.ReadOneBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	beacon := &Beacon{ObjectId: req.ObjectId}
	rsp.Beacon = &device.Beacon{}

	if err := o.Read(beacon, "object_id"); err != nil {
		return err
	}

	rsp.Beacon.Id = beacon.Id
	rsp.Beacon.ObjectId = beacon.ObjectId
	rsp.Beacon.UserId = beacon.UserId
	rsp.Beacon.Mac = beacon.Mac
	rsp.Beacon.Status = beacon.Status
	rsp.Beacon.Password = beacon.Password
	rsp.Beacon.ChangePassword = beacon.ChangePassword
	rsp.Beacon.Type = beacon.Type
	rsp.Beacon.CurrentVersion = beacon.CurrentVersion
	rsp.Beacon.LastSettingVersion = beacon.LastSettingVersion
	rsp.Beacon.ApplyStatus = beacon.ApplyStatus
	rsp.Beacon.CreateAt = beacon.CreateAt
	rsp.Beacon.UpdateAt = beacon.UpdateAt

	return nil
}

func (b *Beacon) ReadByMac(req *device.ReadOneBeaconRequest, rsp *device.ReadOneBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	beacon := &Beacon{Mac: req.Mac}
	rsp.Beacon = &device.Beacon{}

	if err := o.Read(beacon, "mac"); err != nil {
		return err
	}
	rsp.Beacon.Id = beacon.Id
	rsp.Beacon.ObjectId = beacon.ObjectId
	rsp.Beacon.UserId = beacon.UserId
	rsp.Beacon.Mac = beacon.Mac
	rsp.Beacon.Status = beacon.Status
	rsp.Beacon.Password = beacon.Password
	rsp.Beacon.ChangePassword = beacon.ChangePassword
	rsp.Beacon.Type = beacon.Type
	rsp.Beacon.CurrentVersion = beacon.CurrentVersion
	rsp.Beacon.LastSettingVersion = beacon.LastSettingVersion
	rsp.Beacon.ApplyStatus = beacon.ApplyStatus
	rsp.Beacon.CreateAt = beacon.CreateAt
	rsp.Beacon.UpdateAt = beacon.UpdateAt

	return nil
}

func (b *Beacon) ReadByObjectIdAndMac(req *device.ReadOneBeaconRequest, rsp *device.ReadOneBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	beacon := &Beacon{
		Mac:      req.Mac,
		ObjectId: req.ObjectId,
	}
	rsp.Beacon = &device.Beacon{}

	if err := o.Read(beacon, "mac", "object_id"); err != nil {
		return err
	}

	rsp.Beacon.Id = beacon.Id
	rsp.Beacon.ObjectId = beacon.ObjectId
	rsp.Beacon.UserId = beacon.UserId
	rsp.Beacon.Mac = beacon.Mac
	rsp.Beacon.Status = beacon.Status
	rsp.Beacon.Password = beacon.Password
	rsp.Beacon.ChangePassword = beacon.ChangePassword
	rsp.Beacon.Type = beacon.Type
	rsp.Beacon.CurrentVersion = beacon.CurrentVersion
	rsp.Beacon.LastSettingVersion = beacon.LastSettingVersion
	rsp.Beacon.ApplyStatus = beacon.ApplyStatus
	rsp.Beacon.CreateAt = beacon.CreateAt
	rsp.Beacon.UpdateAt = beacon.UpdateAt

	return nil
}

/* 添加Beacon */
func (b *Beacon) Add(req *device.AddBeaconRequest, rsp *device.AddBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	beacon := &Beacon{}

	beacon.ObjectId = req.Beacon.ObjectId
	beacon.UserId = req.Beacon.UserId
	beacon.Mac = req.Beacon.Mac
	beacon.Status = req.Beacon.Status
	beacon.Password = req.Beacon.Password
	beacon.ChangePassword = req.Beacon.ChangePassword
	beacon.Type = req.Beacon.Type
	beacon.CurrentVersion = req.Beacon.CurrentVersion
	beacon.LastSettingVersion = req.Beacon.LastSettingVersion
	beacon.ApplyStatus = req.Beacon.ApplyStatus
	beacon.CreateAt = req.Beacon.CreateAt
	beacon.UpdateAt = req.Beacon.UpdateAt

	if created, id, err := o.ReadOrCreate(beacon, "object_id"); err == nil {
		if created {
			rsp.Result.Code = 200
			rsp.Result.Detail = "create '" + beacon.ObjectId + "' success. 创建'" + beacon.ObjectId + "'成功"
			rsp.Result.Status = "ok"

			return nil
		} else {
			rsp.Result.Code = 201
			rsp.Result.Detail = "create '" + beacon.ObjectId + "' fail. 创建'" + beacon.ObjectId + "'失败"
			rsp.Result.Status = "err"

			return errors.New("AddBeacon", "create '"+beacon.ObjectId+"' fail. 创建'"+beacon.ObjectId+"'失败", 201)
		}
	} else {
		rsp.Result.Code = 202
		rsp.Result.Detail = req.Beacon.ObjectId + "Already exist. ('" + beacon.ObjectId + "'已经存在) ,对应的自增ID为" + strconv.FormatInt(id, 10)
		rsp.Result.Status = "err"

		return err
	}
}

/* 删除Beacon */
func (b *Beacon) Delete(req *device.DeleteBeaconRequest, rsp *device.DeleteBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	beacon := &Beacon{
		ObjectId: req.ObjectId,
		Mac:      req.Mac,
	}

	if err := o.Read(beacon, "object_id", "mac"); err != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "The '" + req.ObjectId + "' Beacon Non-existent."
		rsp.Result.Status = "err"
		return err
	}
	err := o.Begin()
	if err != nil {
		rsp.Result.Code = 202
		rsp.Result.Detail = "Delete Beacon fail."
		rsp.Result.Status = "err"
		return err
	}
	if num, err := o.Delete(beacon, "object_id"); err != nil {
		rsp.Result.Code = 203
		rsp.Result.Detail = "Delete Beacon fail."
		rsp.Result.Status = "err"
		return err
	} else {
		if num != 1 {
			err = o.Rollback()
			rsp.Result.Code = 204
			rsp.Result.Detail = "Delete Beacon fail. Rollback"
			rsp.Result.Status = "err"
			return err
		} else {
			err = o.Commit()
			rsp.Result.Code = 200
			rsp.Result.Detail = "Delete Beacon Success"
			rsp.Result.Status = "ok"
			return err
		}
	}
}

/* 分页读取Beacon */
func (b *Beacon) ReadPaging(req *device.ReadPagingBeaconRequest, rsp *device.ReadPagingBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	qs := o.QueryTable("beacon")

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
		}
	}

	if req.Limit > 0 && req.Offset >= 0 {
		qs = qs.Limit(req.Limit, (req.Offset-1)*req.Limit)
	} else {
		qs = qs.Limit(10, 0)
	}

	var beacons []*Beacon
	qs.All(&beacons)

	for i := 0; i < len(beacons); i++ {
		bc := &device.Beacon{}

		bc.Id = beacons[i].Id
		bc.ObjectId = beacons[i].ObjectId
		bc.UserId = beacons[i].UserId
		bc.Mac = beacons[i].Mac
		bc.Status = beacons[i].Status
		bc.Password = beacons[i].Password
		bc.ChangePassword = beacons[i].ChangePassword
		bc.Type = beacons[i].Type
		bc.CurrentVersion = beacons[i].CurrentVersion
		bc.LastSettingVersion = beacons[i].LastSettingVersion
		bc.ApplyStatus = beacons[i].ApplyStatus
		bc.CreateAt = beacons[i].CreateAt
		bc.UpdateAt = beacons[i].UpdateAt

		rsp.Beacons = append(rsp.Beacons, bc)
	}
	if cnt, ce := o.QueryTable("beacon").Count(); ce == nil {
		rsp.Total = cnt
	}

	return nil
}

/* 更新Beacon */
func (b *Beacon) Update(req *device.UpdateBeaconRequest, rsp *device.UpdateBeaconResponse) error {
	UseingDB(req.DatabaseUrl)
	beacon := &Beacon{}

	beacon.ObjectId = req.Beacon.ObjectId
	beacon.UserId = req.Beacon.UserId
	beacon.Mac = req.Beacon.Mac
	beacon.Status = req.Beacon.Status
	beacon.Password = req.Beacon.Password
	beacon.ChangePassword = req.Beacon.ChangePassword
	beacon.Type = req.Beacon.Type
	beacon.CurrentVersion = req.Beacon.CurrentVersion
	beacon.LastSettingVersion = req.Beacon.LastSettingVersion
	beacon.ApplyStatus = req.Beacon.ApplyStatus
	beacon.CreateAt = req.Beacon.CreateAt
	beacon.UpdateAt = req.Beacon.UpdateAt

	if err := o.Read(beacon, "object_id"); err != nil {
		rsp.Result.Detail = "Beacon Non-existent."
		rsp.Result.Code = 201
		rsp.Result.Status = "ok"
		return nil
	}
	rsp.Result = &device.Result{Id: beacon.ObjectId}

	if _, err := o.Update(beacon); err != nil {
		rsp.Result.Detail = "Update Beacon fail."
		rsp.Result.Code = 500
		rsp.Result.Status = "err"
		return err
	}
	rsp.Result.Detail = "Update Beacon Success."
	rsp.Result.Code = 200
	rsp.Result.Status = "ok"
	return nil

}

func (bs *Beacon) UpdateApplyStatus(req *device.UpdateBeaconApplyStatusRequest, rsp *device.UpdateBeaconApplyStatusResponse) error {
	UseingDB(req.DatabaseUrl)
	beacon := &Beacon{
		ObjectId: req.ObjectId,
		Mac:      req.Mac,
	}

	if err := o.Read(beacon, "object_id", "mac"); err != nil {
		rsp.Result.Code = 201
		rsp.Result.Detail = "The '" + req.ObjectId + "' UpdateBeaconApplyStatus Non-existent."
		rsp.Result.Status = "err"
		return err
	}

	beacon.ApplyStatus = 1
	if _, err := o.Update(beacon); err != nil {
		rsp.Result.Code = 202
		rsp.Result.Detail = "Update Update UpdateBeaconApplyStatus Apply Version Status fail."
		rsp.Result.Status = "err"
		return err
	} else {
		rsp.Result.Code = 200
		rsp.Result.Detail = "Update Update UpdateBeaconApplyStatus Apply Version Status Success."
		rsp.Result.Status = "ok"
		return nil
	}
}

func (b *Beacon) ReadOneFilter(objectId string) *device.Beacon {
	beacon := &device.Beacon{}

	err := o.QueryTable(beacon).Filter("id", 1).One(beacon) // 返回 QuerySeter

	if err == orm.ErrMultiRows {
		// 多条的时候报错
		fmt.Printf("多条的时候报错")
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		fmt.Printf("没有找到记录")
	}
	return beacon
}
