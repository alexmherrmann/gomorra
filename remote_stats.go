package gomorra

import "strings"

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
	return count-1, nil
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

func (r *Remote) GetLoadPercentage() (float32, error) {
	panic("implement me")
}