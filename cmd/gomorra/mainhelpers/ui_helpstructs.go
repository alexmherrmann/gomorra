package mainhelpers

import (
	"github.com/alexmherrmann/gomorra"
	"strings"
	"sort"
)

type DisplayableStat struct {
	Config      gomorra.HostConfig
	Remote      *gomorra.Remote
	LoadChannel chan gomorra.StatResult
}

type NamedPercentageResult struct {
	Name   string
	Result int
	Err    error
}

type nicety struct {
	hostPrettyName   string
	loadPercentage   int
	availableMemInGb float32
	err              error
}

func (a statStore) copy() statStore {
	newStore := make(statStore)
	state.mapMutex.Lock()
	for key := range a {
		newStore[key] = a[key]
	}
	state.mapMutex.Unlock()

	return newStore
}

func (dontuse statStore) getNicely() []nicety {

	copied := dontuse.copy()

	toReturn := make([]nicety, 0, 5)
	helperMap := make(map[string]*nicety)


	for key := range copied {
		val := copied[key]
		split := strings.Split(key, "\t")
		//state.logger.Printf("splits: %#v", split)
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
			returnNicety.loadPercentage = val.Result
			break
		case NicetyAvailable:
			returnNicety.availableMemInGb = float32(val.Result) / 1024. / 1024.
			break
		case NicetyError:
			returnNicety.err = val.Err
			break
		default:
			state.logger.Println("Got a weird type: ", theType)
		}

		//toReturn = append(toReturn, *returnNicety)
	}

	keys := make([]string, 0, 5)
	for key := range helperMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := helperMap[key]
		toReturn = append(toReturn, *val)
	}

	state.logger.Println("Returning: ", toReturn)
	return toReturn
}
