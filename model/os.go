package model

type (
	OS string
)

const (
	Android OS = "android"
	iOS     OS = "ios"
	Windows OS = "windows"
	Mac     OS = "mac"
	Linux   OS = "linux"
)

var OperationalSystems = map[string]OS{
	"android": Android,
	"ios":     iOS,
	"windows": Windows,
	"mac":     Mac,
	"linux":   Linux,
}
