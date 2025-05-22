package model

type (
	Device string
)

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
