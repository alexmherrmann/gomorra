package mainhelpers

import (
	t "github.com/gizak/termui"
)

type NamedPercentageResult struct {
	Name   string
	Result int
}

func BuildStatGrid(toBuildFrom []NamedPercentageResult) *t.Grid {
	grid := t.NewGrid()

	for _, result := range toBuildFrom {
		gauge := t.NewGauge()
		gauge.BorderLabel = result.Name
		gauge.Percent = result.Result

		grid.AddRows(
			t.NewRow(
				t.NewCol(12, 0, gauge),
			),
		)
	}

	return grid

}

func ShowStats(toDisplay []NamedPercentageResult) {
	err := t.Init()
	if err != nil {
		panic(err.Error())
	}

	grid := BuildStatGrid(toDisplay)
	t.Render(grid)
}
