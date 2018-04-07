package mainhelpers

import (
	t "github.com/gizak/termui"
	"sync"
	"log"
)

var state = struct {
	mutex       sync.Mutex
	percentages map[string]int
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
		state.logger.Printf("Received result %s = %d\n", received.Name, received.Result)
		// TODO: rethink having to register
		state.mutex.Lock()
		state.percentages[received.Name] = received.Result
		state.mutex.Unlock()
	}
}

func buildGridFromCurrentValues() {
	grid := t.NewGrid()
	state.logger.Println("Building grid from values")

	grid.Width = t.TermWidth()
	grid.X = 0
	grid.Y = 0
	grid.BgColor = t.ColorCyan

	//rows := make([]*t.Row, len(state.percentages))

	for key := range state.percentages {
		state.logger.Println("have key: " + key)
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
