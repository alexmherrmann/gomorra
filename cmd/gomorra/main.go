package main

import (
	g "github.com/alexmherrmann/gomorra"
	. "github.com/alexmherrmann/gomorra/mainhelpers"
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
	logfile, err := os.OpenFile("gomorra.log", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic("couldn't open log file")
	}
	mainLogger := log.New(logfile, "gomorra ", log.LstdFlags)
	mainLogger.Println("Gomorra started")
	return err, mainLogger
}

func main() {

	err, mainLogger := getMainLogger()

	//if runtime.GOOS == "windows" {
	//	fmt.Println("SIGNAL HANDLING PROBABLY WON'T WORK")
	//} else {
	//	go interruptListener(mainLogger)
	//	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL)
	//}

	var configFilePath string
	flag.StringVar(&configFilePath, "file", "config.json", "The path to the configuration json file")

	config, err := g.ReadConfigFile(configFilePath)

	if err != nil {
		mainLogger.Fatalln("Couldn't open config file: " + configFilePath)
	}

	mainLogger.Printf("Have %d configs\n", len(config.Hosts))
	loadChannel := make(chan g.StatResult)

	firstRemote, err := g.GetRemoteFromHostConfig(config.Hosts[0])
	if err != nil {
		mainLogger.Fatalln(err.Error())
	}

	err = firstRemote.Open()
	if err != nil {
		mainLogger.Fatalln(err.Error())
	}

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

	go t.Loop()

	exit := false;
	// TODO: Make this accept multiple
	for true {

		go firstRemote.GetLoadMinuteAvg(loadChannel)
		result := <-loadChannel
		timeString := time.Now().Format("3:04:05 pm")

		if load, ok := g.CheckFloat(result); ok {
			mainLogger.Printf("%s: %.3f\n", timeString, load)
			// TODO: This is ugly pls change
			ShowStats([]NamedPercentageResult{NamedPercentageResult{config.Hosts[0].Prettyname, int(load * 100)}})
		} else {
			mainLogger.Printf("%s: couldn't get load\n", timeString)
		}

		select {
		case <-quitChan:
			//close(quitChan)
			exit = true
			break
		case <-time.After(5 * time.Second):
			continue
		}

		if exit {
			break
		}


	}

	mainLogger.Println("shutting down")
}
