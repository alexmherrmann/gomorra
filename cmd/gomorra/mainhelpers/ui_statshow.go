package mainhelpers

import (
	t "github.com/gizak/termui"
	"sync"
	"log"
	"sort"
	"fmt"
)

type statStore map[string]NamedPercentageResult

var state = struct {
	mapMutex    sync.Mutex
	percentages statStore
	errors      map[string]error
	grid        *t.Grid
	logger      *log.Logger
}{}

func init() {
	state.percentages = make(statStore)
}

const NicetyLoadAvg = "load avg"
const NicetyAvailable = "available"
const NicetyError = "error"

func BeginListen(c <-chan NamedPercentageResult, logger *log.Logger) {
	go ReceiveResults(c)
	err := t.Init()
	if err != nil {
		panic(err.Error())
	}
	go t.Loop()
	state.logger = logger
}

func ReceiveResults(listenerChannel <-chan NamedPercentageResult) {
	for received := range listenerChannel {
		state.logger.Println("Received result: ", received)
		state.mapMutex.Lock()
		state.percentages[received.Name] = received
		state.mapMutex.Unlock()
	}
}

func buildGridFromCurrentValues() {
	grid := t.NewGrid()

	grid.Width = t.TermWidth()
	grid.X = 0
	grid.Y = 0
	grid.BgColor = t.ColorCyan

	//rows := make([]*t.Row, len(state.percentages))

	state.logger.Printf("have percentages: %d\n", len(state.percentages))
	var keys []string = make([]string, 0, 5)
	for key := range state.percentages {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, nicety := range state.percentages.getNicely() {

		if nicety.err == nil {
			loadGauge := t.NewGauge()
			memstr := fmt.Sprintf("%3.2f GB", nicety.availableMemInGb);
			availableMemory := t.NewPar(memstr)

			availableMemory.BorderLabel = nicety.hostPrettyName + " available mem"
			availableMemory.Height = 3

			loadGauge.BorderLabel = nicety.hostPrettyName + " load"
			loadGauge.Percent = nicety.loadPercentage
			loadGauge.Height = 3

			grid.AddRows(
				t.NewRow(
					t.NewCol(9, 0, loadGauge),
					t.NewCol(3, 0, availableMemory),
				),
			)
		} else {
			errorPar := t.NewPar(nicety.err.Error())
			errorPar.Height = 3
			errorPar.BorderLabel = nicety.hostPrettyName + " error"

			grid.AddRows(t.NewRow(
				t.NewCol(12, 0, errorPar),
			))
		}
	}

	grid.Align()
	state.grid = grid
}

func ShowStats() {
	buildGridFromCurrentValues()
	//state.logger.Println("rendering values!");
	t.Render(state.grid)
}
