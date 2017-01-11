package apiEntity

import (
	"beaconCloud/rpcServer/device-srv/proto/device"
	"github.com/micro/go-micro/client"
)

var (
	DatabaseUrl = ""
	cl device.DeviceClient
)

func init() {
	cl = device.NewDeviceClient("go.micro.srv.device-srv", client.DefaultClient)
}