package main

import (
	g "github.com/alexmherrmann/gomorra"
	. "github.com/alexmherrmann/gomorra/cmd/gomorra/mainhelpers"
	"flag"
	"log"
	"os"
	"time"
	t "github.com/gizak/termui"
	"sync"
)

var signalChan = make(chan os.Signal)

func interruptListener(logger *log.Logger) {
	<-signalChan
	logger.Println("received interrupt")
	os.Exit(0)
}

func getMainLogger() (error, *log.Logger) {
	os.Remove("gomorra.log")

	logfile, err := os.OpenFile("gomorra.log", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic("couldn't open log file")
	}
	mainLogger := log.New(logfile, "gomorra ", log.LstdFlags)
	mainLogger.Println("Gomorra started")
	return err, mainLogger
}

func GetAllDisplayableStats(configFilePath string) []DisplayableStat {
	readConfig, err := g.ReadConfigFile(configFilePath)

	if err != nil {
		mainLogger.Fatalln("Couldn't open readConfig file: " + configFilePath)
	}
	listToReturn := make([]DisplayableStat, 0)

	for _, config := range readConfig.Hosts {
		remote, err := g.GetRemoteFromHostConfig(config)
		if err != nil {
			mainLogger.Fatalln(err.Error())
		}

		err = remote.Open()
		if err != nil {
			mainLogger.Fatalln(err.Error())
		}

		listToReturn = append(listToReturn, DisplayableStat{
			Remote:      remote,
			Config:      config,
			LoadChannel: make(chan g.StatResult),
		})
	}

	return listToReturn
}

var mainLogger log.Logger

func main() {

	err, mainLogger := getMainLogger()
	if err != nil {
		mainLogger.Fatalln(err.Error())
	}

	// TODO: Re-enable this at some point
	//if runtime.GOOS == "windows" {
	//	fmt.Println("SIGNAL HANDLING PROBABLY WON'T WORK")
	//} else {
	//	go interruptListener(mainLogger)
	//	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL)
	//}

	var configFilePath string
	flag.StringVar(&configFilePath, "file", "config.json", "The path to the configuration json file")
	statStuff := GetAllDisplayableStats(configFilePath)


	mainLogger.Printf("Have %d configs\n", len(statStuff))

	quitChan := make(chan interface{})
	quitOnceChan := sync.Once{}

	stopped := false

	t.Handle("/sys/kbd/q", func(t.Event) {
		// press q to quit
		quitOnceChan.Do(func() {
			if !stopped {
				stopped = true
				mainLogger.Println("Received quit")
				quitChan <- true
				t.StopLoop()
			}
		})
	})

	// I dunno, 3 seems good
	displayResultListenerChannel := make(chan NamedPercentageResult, len(statStuff))
	go BeginListen(displayResultListenerChannel, mainLogger)
	exit := false;

	go func() {
		for range time.Tick(2 * time.Second) {
			ShowStats()
		}
	}()

	mainLogger.Println("beginning listen")
	for true {

		for _, toDisplay := range statStuff {
			remote := toDisplay.Remote
			config := toDisplay.Config
			channel := toDisplay.LoadChannel

			go remote.GetLoadMinuteAvg(channel)
			result := <-channel

			floatVal, ok := g.CheckFloat(result)
			if ok {
				mainLogger.Printf("Got a result for %s! %.3f\n", config.Prettyname, floatVal)
				namedResult := NamedPercentageResult{
					Name: config.Prettyname,
					Result: int(floatVal * 100),
				}
				displayResultListenerChannel <- namedResult


			} else {
				mainLogger.Println("Got a bad result from the channel")
			}
		}

		select {
		case <-quitChan:
			exit = true
			break
		case <-time.After(5 * time.Second):
			continue
		}

		if exit {
			break
		}
		ShowStats()

	}

	mainLogger.Println("shutting down")
}
