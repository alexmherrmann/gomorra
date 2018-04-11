package mainhelpers

import (
	"github.com/alexmherrmann/gomorra"
	"strings"
)

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

func (a statStore) getNicely() []nicety {
	toReturn := make([]nicety, 0, 5)
	helperMap := make(map[string]*nicety)

	for key := range a {
		val := a[key]
		split := strings.Split(key, "\t")
		state.logger.Printf("splits: %#v", split)
		name := split[0]
		theType := split[1]

		var returnNicety *nicety
		var exist bool

		if returnNicety, exist = helperMap[name]; !exist {
			returnNicety = new(nicety)
			helperMap[name] = returnNicety
			returnNicety.hostPrettyName = name
		}


		switch theType {
		case NicetyLoadAvg:
			returnNicety.loadPercentage = val
			break
		case NicetyAvailable:
			returnNicety.availableMemInGb = float32(val) / 1024. / 1024.
			break
		}

		//toReturn = append(toReturn, *returnNicety)
	}

	for key := range helperMap {
		val := helperMap[key]
		toReturn = append(toReturn, *val)
	}

	state.logger.Println("Returning: ", toReturn)
	return toReturn
}
