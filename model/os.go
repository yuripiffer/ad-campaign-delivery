package model

type (
	OS string
)

// REMINDER: also insert the OS in map OperationalSystems whenever
// a new OS is added as a constant.
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
