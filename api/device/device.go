package device

import (
	"beaconCloud/api/device/apiEntity"
)

type Device struct {
	apiEntity.Gateway
	apiEntity.Beacon
	apiEntity.BeaconSetting
}
