package mainhelpers

import "github.com/alexmherrmann/gomorra"

type DisplayableStat struct {
	Config gomorra.HostConfig
	Remote *gomorra.Remote
	LoadChannel chan gomorra.StatResult
}

type NamedPercentageResult struct {
	Name   string
	Result int
}
