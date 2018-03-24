package main

import (
	"fmt"
	_ "github.com/gizak/termui"
	gom "alexmherrmann.com/gomorra"
	"math/rand"
	"log"
)

//TODO: make this channels or something cool
func randomFloat() float32 {
	return rand.Float32()
}

// TODO: make these more type safe
func printStats(gettable gom.ComputerStatGettable) {
	resultChan := make(chan gom.StatResult)

	go gettable.GetCores(resultChan)

	result := <-resultChan
	gom.FatalErr(result.Err)

	cores := result.GenericResult.(int)
	fmt.Printf("Have %d cores\n", cores)


	go gettable.GetLoadMinuteAvg(resultChan)
	result = <-resultChan
	gom.FatalErr(result.Err)
	fmt.Printf("Our last minute load percentage: %%%.2f\n", result.GenericResult.(float32) * 100)

}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	localhost := gom.Remote{
		Hostname: "127.0.0.1:22",
	}

	err := localhost.Open("alex", "/home/alex/.ssh/id_rsa")
	gom.FatalErr(err)

	printStats(&localhost)
}
