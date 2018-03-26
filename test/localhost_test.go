package test

import "testing"

import (
	gom "alexmherrmann.com/gomorra"
)

// TODO: make these more type safe
func printStats(gettable gom.ComputerStatGettable, t *testing.T) {
	resultChan := make(chan gom.StatResult)

	failErr := func(err error) {
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	go gettable.GetCores(resultChan)

	result := <-resultChan
	failErr(result.Err)

	cores, ok := gom.CheckInt(result)
	if !ok {
		t.Fatal("didn't get int")
	}
	t.Logf("Have %d cores\n", cores)

	go gettable.GetLoadMinuteAvg(resultChan)
	result = <-resultChan
	failErr(result.Err)
	t.Logf("Our last minute load percentage: %%%.2f\n", result.GenericResult.(float32)*100)

	go gettable.GetTotalMemory(resultChan)
	result = <-resultChan
	failErr(result.Err)

	var totalKb int
	totalKb, ok = gom.CheckInt(result)

	if !ok {
		t.Fatal("Didn't get an int back")
	}
	t.Logf("Have %d Mb of total memory\n", totalKb/1024)

	go gettable.GetAvailableMemory(resultChan)
	result = <-resultChan
	failErr(result.Err)

	var freeKb int
	freeKb, ok = gom.CheckInt(result)
	if !ok {
		t.Fatal("Didn't get an int back")
	}
	t.Logf("Have %d Mb of available memory\n", freeKb/1024)

	go gom.GetUsedMemory(resultChan, gettable)
	result = <-resultChan
	failErr(result.Err)

	available, ok := gom.CheckInt(result)
	if !ok {
		t.Fatal("Didn't get an int back")
	}
	t.Logf("Have %d Mb of used memory", available/1024)
}

// This only works on my local machine
func TestLocalhost(t *testing.T) {
	localhost := gom.Remote{
		Hostname: "localhost:22",
	}

	err := localhost.Open("alex", "/home/alex/.ssh/id_rsa")
	if err != nil {
		t.Error("Couldn't successfully open localhost:\n" + err.Error())
	}

	printStats(&localhost, t)

}
