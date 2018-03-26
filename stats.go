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
	msg    string
}

func (i IncorrectTypeError) Error() string {
	if len(i.msg) == 0 {
		return fmt.Sprintf("We wanted a type of: %s", i.wanted)
	} else {
		return fmt.Sprintf("Wanted %s: %s", i.wanted, i.msg)
	}
}

func Wanted(wanted string) IncorrectTypeError {
	return IncorrectTypeError{wanted, ""}
}

func WantedMessage(wanted string, msg string) IncorrectTypeError {
	return IncorrectTypeError{wanted, msg}
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
	// Get the amount of available memory
	GetAvailableMemory(chan StatResult)
}

func GetUsedMemory(channel chan StatResult, gettable ComputerStatGettable) {
	availableChan := make(chan StatResult)
	totalChan := make(chan StatResult)

	go gettable.GetTotalMemory(totalChan)
	go gettable.GetAvailableMemory(availableChan)

	availableResult := <-availableChan
	totalResult := <-totalChan

	if availableResult.Err != nil {
		channel <- availableResult
		return
	}

	if totalResult.Err != nil {
		channel <- totalResult
		return
	}

	total, ok := CheckInt(totalResult)
	CheckInt(totalResult)
	if !ok {
		channel <- StatResult{Err: WantedMessage("int", "total")}
		return
	}

	available, ok := CheckInt(availableResult)

	if !ok {
		channel <- StatResult{Err: WantedMessage("int", "available")}
		return
	}

	channel <- StatResult{GenericResult: total - available}
}
