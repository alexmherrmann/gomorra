package gomorra

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var NotImplementedErr = errors.New("Not implemented")

func errResult(err error) StatResult {
	return StatResult{
		Err:err,
		Type: TError,
	}
}

// internal function that forces ssh connection
func (r *Remote) getCores() (int, error) {
	coreData, err := r.readFileFromSystem("/proc/stat")
	if err != nil {
		return -1, err
	}

	splitStrings := strings.Split(coreData.String(), "\n")

	var count int = 0
	for _, line := range splitStrings {
		if strings.Index(line, "cpu") >= 0 {
			count++
		} else {
			break
		}
	}

	// we take one off because the first cpu line is bad
	return count - 1, nil
}

// returns the 3 load percentages for 1, 5 and 15 minutes
func (r *Remote) getLoads() ([3]float32, error) {
	var badReturn = [3]float32{1, 1, 1}
	//TODO: implement error handling
	newSesh, err := r.client.NewSession()
	if err != nil {
		return badReturn, err
	}

	readBytes := new(bytes.Buffer)
	newSesh.Stdout = readBytes

	//TODO: change this to use the readFileFromSystem function
	//TODO: implement error handling
	err = newSesh.Run("/usr/bin/env cat /proc/loadavg")
	if err != nil {
		return badReturn, err
	}

	var avg1, avg2, avg3 float32

	fmt.Sscanf(readBytes.String(), "%f %f %f", &avg1, &avg2, &avg3)

	return [3]float32{avg1, avg2, avg3}, nil
}

// This will only go to the server to get the number of cores if it hasn't already
// Statresult will have an int
func (r *Remote) GetCores(channel chan<- StatResult) {
	toReturn := StatResult {
		Type: TCores,
	}

	if r.cores == nil {
		cores, err := r.getCores()
		if err == nil {
			r.cores = new(int)
			*r.cores = cores
			toReturn.GenericResult = *r.cores
			channel <- toReturn
			return
		}
		channel <- errResult(err)
		return
	} else {
		toReturn.GenericResult = *r.cores
		channel <- toReturn
		return
	}
}

// Get the last minutes load percentage as FLOAT
func (r *Remote) GetLoadMinuteAvg(channel chan<- StatResult) {

	coreResult := make(chan StatResult)
	go r.GetCores(coreResult)
	result := <-coreResult

	var cores int
	if v, ok := CheckInt(result); ok {
		cores = v
	} else {
		channel <- StatResult{Err: Wanted("int")}
		return
	}

	if result.Err != nil {
		channel <- StatResult{Err: result.Err}
		return
	}

	avgs, err := r.getLoads()
	if err != nil {
		channel <- StatResult{Err: err}
		return
	}

	channel <- StatResult{
		GenericResult: avgs[0] / float32(cores),
		Type:          TLoadMinuteAvg,
	}
}

func (r *Remote) getMeminfo() (string, error) {
	meminfo, err := r.readFileFromSystem("/proc/meminfo")
	if err != nil {
		return "", err
	}

	return meminfo.String(), nil
}

const stringFormat = `MemTotal: %d kB
MemFree: %d kB
MemAvailable: %d kB`

// Returns the total memory in kb
func (r *Remote) getTotalMemory() (int, error) {
	memInfoString, err := r.getMeminfo()

	if err != nil {
		return 0, err
	}

	var total int

	fmt.Sscanf(memInfoString, stringFormat, &total)

	return total, nil
}

// Will contain INT of Total kb
func (r *Remote) GetTotalMemory(channel chan<- StatResult) {
	// go and fetch the total amount of memory
	if r.totalMemKb == nil {
		totalMem, err := r.getTotalMemory()
		if err != nil {
			channel <- StatResult{Err: err}
			return
		}
		r.totalMemKb = new(int)
		*r.totalMemKb = totalMem
		channel <- StatResult{
			GenericResult: *r.totalMemKb,
			Type:TTotalMemory,
		}
		return
	}

	channel <- StatResult{
		GenericResult: *r.totalMemKb,
		Type:TTotalMemory,
	}

}

// Statresult will have INT of Free KB
func (r *Remote) GetFreeMemory(channel chan<- StatResult) {
	memInfoString, err := r.getMeminfo()

	if err != nil {
		channel <- StatResult{Err: err}
		return
	}

	var free int

	// we use free twice to overwrite the first one so we don't need to have a fake variable
	fmt.Sscanf(memInfoString, stringFormat, &free, &free)

	channel <- StatResult{
		GenericResult: free,
		Type:          TFreeMemory,
	}
}

// Statresult should contain an INT of the amount of available kilobytes
func (r *Remote) GetAvailableMemory(channel chan<- StatResult) {
	memInfoString, err := r.getMeminfo()

	if err != nil {
		channel <- StatResult{Err: err}
		return
	}

	var available int

	// we use available thrice to overwrite the first one so we don't need to have a fake variable
	fmt.Sscanf(memInfoString, stringFormat, &available, &available, &available)

	channel <- StatResult{
		GenericResult: available,
		Type:          TAvailableMemory,
	}
}

func (r *Remote) GetNumberOfSecondsUptime(channel chan<- StatResult) {
	channel <- StatResult{Err:NotImplementedErr}
}
