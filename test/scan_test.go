package test

import (
	"testing"
	"fmt"
	"strconv"
)

func TestMemScan(t *testing.T) {
	const exampleString = `MemTotal:       16289440 kB
MemFree:        10104424 kB
MemAvailable:   13019388 kB
Buffers:          143008 kB
Cached:          3174976 kB
SwapCached:            0 kB
Active:          2878996 kB`
	const stringFormat = `MemTotal: %d kB
MemFree: %d kB
MemAvailable: %d kB`

	var total, free, available int
	fmt.Sscanf(exampleString, stringFormat, &total, &free, &available)

	if total != 16289440 {
		t.Error("total not read correctly: " + strconv.Itoa(total))
	}

	if free != 10104424 {
		t.Error("free not read correctly: " + strconv.Itoa(free))
	}

	if available != 13019388 {
		t.Error("available not read correctly: " + strconv.Itoa(available))
	}
}
