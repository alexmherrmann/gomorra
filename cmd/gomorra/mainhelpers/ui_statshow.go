package mainhelpers

import (
	t "github.com/gizak/termui"
	"sync"
	"log"
	"sort"
)

type statStore map[string]int
var state = struct {
	mutex       sync.Mutex
	percentages statStore
	grid        *t.Grid
	logger      *log.Logger
}{}

func Register(name string) {
	state.percentages[name] = 0
}

func init() {
	state.percentages = make(map[string]int)
}

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
		// TODO: rethink having to register
		state.mutex.Lock()
		state.percentages[received.Name] = received.Result
		state.mutex.Unlock()
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

	for _, key := range keys {
		gauge := t.NewGauge()
		value := state.percentages[key]
		gauge.BorderLabel = key
		gauge.Percent = value
		gauge.Height = 3
		grid.AddRows(
			t.NewRow(
				t.NewCol(12, 0, gauge),
			),
		)
	}

	grid.Align()
	state.grid = grid
}

func ShowStats() {
	buildGridFromCurrentValues()
	//state.logger.Println("rendering values!");
	t.Render(state.grid)
}
