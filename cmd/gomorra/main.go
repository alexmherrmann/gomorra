package main

import (
	"fmt"
	_ "github.com/gizak/termui"
	gom "alexmherrmann.com/gomorra"
	"math/rand"
)

//TODO: make this channels or something cool
func randomFloat() float32 {
	return rand.Float32()
}

func printStats(gettable gom.ComputerStatGettable) {
	cores, err := gettable.GetCores()
	gom.FatalErr(err)
	fmt.Printf("Have %d cores\n", cores)

	
}

func main() {

	localhost := gom.Remote{
		Hostname: "127.0.0.1:22",
	}

	err := localhost.Open("alex", "/home/alex/.ssh/id_rsa")
	gom.FatalErr(err)

	printStats(&localhost)
}
