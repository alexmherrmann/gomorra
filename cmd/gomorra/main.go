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

func GetAllDisplayableStats(configFilePath string, channel chan <- NamedPercentageResult) []DisplayableStat {
	readConfig, err := g.ReadConfigFile(configFilePath)

	if err != nil {
		mainLogger.Fatalln("Couldn't open readConfig file: " + configFilePath)
	}
	listToReturn := make([]DisplayableStat, 0)

	for _, config := range readConfig.Hosts {
		remote, err := g.GetRemoteFromHostConfig(config, mainLogger)
		if err != nil {
			mainLogger.Fatalln(err.Error())
		}

		err = remote.Open()
		if err != nil {
			channel <- NamedPercentageResult{Name:config.Prettyname + "\t" + NicetyError, Err: err}
			mainLogger.Println("Got error on ", config.Prettyname ,":\n", err.Error())
		}

		listToReturn = append(listToReturn, DisplayableStat{
			Remote:      remote,
			Config:      config,
			LoadChannel: make(chan g.StatResult),
		})
	}

	return listToReturn
}

var mainLogger *log.Logger

func DoGetStandardLoad(listenerChannel chan<- NamedPercentageResult, stat DisplayableStat) {

	waiter := sync.WaitGroup{}
	waiter.Add(2)

	go func(stat DisplayableStat) {
		defer waiter.Done()
		loadAvgChannel := make(chan g.StatResult)
		go stat.Remote.GetLoadMinuteAvg(loadAvgChannel)

		result := <-loadAvgChannel

		if floatVal, ok := g.CheckFloat(result); ok {
			listenerChannel <- NamedPercentageResult{
				Name:   stat.Config.Prettyname + "\t" + NicetyLoadAvg,
				Result: int(floatVal * 100),
			}
		}
	}(stat)

	go func(stat DisplayableStat) {
		defer waiter.Done()
		availableMemChannel := make(chan g.StatResult)
		go stat.Remote.GetAvailableMemory(availableMemChannel)

		result := <-availableMemChannel

		if intVal, ok := g.CheckInt(result); ok {
			listenerChannel <- NamedPercentageResult{
				Name:   stat.Config.Prettyname + "\t" + NicetyAvailable,
				Result: intVal,
			}
		}
	}(stat)

	waiter.Wait()

}

func main() {

	var err error
	err, mainLogger = getMainLogger()
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

	// TODO: the 15 hardcoded in here should be a little smarter
	displayResultListenerChannel := make(chan NamedPercentageResult, 15)
	go BeginListen(displayResultListenerChannel, mainLogger)

	var configFilePath string
	flag.StringVar(&configFilePath, "file", "config.json", "The path to the configuration json file")
	statStuff := GetAllDisplayableStats(configFilePath, displayResultListenerChannel)

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

	exit := false;
	go func() {
		for range time.Tick(2 * time.Second) {
			ShowStats()
		}
	}()

	group := sync.WaitGroup{}

	mainLogger.Println("Beginning listen")
	for true {

		for _, toDisplay := range statStuff {
			group.Add(1)
			//mainLogger.Printf("Starting routine for %s\n", toDisplay.Remote.Hostname)
			go func(d DisplayableStat) {
				defer group.Done()
				DoGetStandardLoad(displayResultListenerChannel, d)
			}(toDisplay)
		}

		group.Wait()

		select {
		case <-quitChan:
			exit = true
			break
		case <-time.After(3 * time.Second):
			continue
		}

		if exit {
			break
		}
	}

	mainLogger.Println("shutting down")
}
