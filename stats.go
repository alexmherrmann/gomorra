package gomorra

import (
	"fmt"
)

const (
	Cores         = iota
	LoadMinuteAvg
	TotalMemory
	FreeMemory
)

type IncorrectTypeError struct {
	wanted string
}

func (i IncorrectTypeError) Error() string {
	return fmt.Sprintf("We wanted a type of: %s", i.wanted)
}

func Wanted(wanted string) IncorrectTypeError {
	return IncorrectTypeError{wanted}
}

type StatResult struct {
	// This is probably a really bad design decision but let's see how it happens
	GenericResult interface{}
	Err           error
}

type ComputerStatGettable interface {
	// Get the number of cores in the system
	GetCores(chan StatResult)
	// Get the last minute load average
	GetLoadMinuteAvg(chan StatResult)
	// Get the total amount of memory on the system
	GetTotalMemory(chan StatResult)
	// Get the amount of free memory
	GetFreeMemory(chan StatResult)
}
