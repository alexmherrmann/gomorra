package mainhelpers

import "github.com/alexmherrmann/gomorra"

type DisplayableStat struct {
	Config      gomorra.HostConfig
	Remote      *gomorra.Remote
	LoadChannel chan gomorra.StatResult
}

type NamedPercentageResult struct {
	Name   string
	Result int
}

type nicety struct {
	hostPrettyName   string
	loadPercentage   int
	availableMemInGb float32
}

//func (a statStore) getNicely() []nicety {
//
//}
