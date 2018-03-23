package gomorra

import (
	"strings"
	"bytes"
	"fmt"
)

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
		}
	}

	// we take one off because the first cpu line is bad
	return count - 1, nil
}

// returns the 3 load percentages for 1, 5 and 15 minutes
func (r *Remote) getLoads() ([3]float32, error) {
	var badReturn = [3]float32{1,1,1}
	//TODO: implement error handling
	newSesh, err := r.client.NewSession()
	if err != nil {
		return badReturn, err
	}

	readBytes := new(bytes.Buffer)
	newSesh.Stdout = readBytes

	//TODO: implement error handling
	err = newSesh.Run("/usr/bin/env cat /proc/loadavg")
	if err != nil {
		return badReturn, err
	}

	//DataLogger.Printf("Got [%s] for loadavg", strings.Trim(readBytes.String(), "\n"))

	var avg1, avg2, avg3 float32

	fmt.Sscanf(readBytes.String(), "%f %f %f", &avg1, &avg2, &avg3)

	return [3]float32{avg1, avg2, avg3}, nil
}

// This will only go to the server to get the number of cores if it hasn't already
func (r *Remote) GetCores() (int, error) {
	if r.cores == nil {
		cores, err := r.getCores()
		if err == nil {
			r.cores = new(int)
			*r.cores = cores
			return *r.cores, nil
		}
		return -1, err
	} else {
		return *r.cores, nil
	}
}

// Get the last minutes load percentage
func (r *Remote) GetLoadAvgPercentage() (float32, error) {
	cores, err := r.GetCores()
	if err != nil {
		return 0, err
	}

	avgs, err := r.getLoads()
	if err != nil {
		return 0, err
	}

	return avgs[0]/float32(cores), nil
}

func (r *Remote) GetTotalMemory() (int, error) {
	panic("implement me")
}

func (r *Remote) GetFreeMemory() (int, error) {
	panic("implement me")
}
