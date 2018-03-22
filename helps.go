package gomorra

import "log"

func FatalErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}