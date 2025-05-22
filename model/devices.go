package model

type (
	Device string
)

// REMINDER: also insert the device in map Devices whenever
// a new device is added as a constant.
const (
	Mobile  Device = "mobile"
	Desktop Device = "desktop"
	Tablet  Device = "tablet"
)

var Devices = map[string]Device{
	"mobile":  Mobile,
	"desktop": Desktop,
	"tablet":  Tablet,
}
