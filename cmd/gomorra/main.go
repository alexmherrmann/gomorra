package main
import (
	"fmt"
	tui "github.com/gizak/termui"
	gom "alexmherrmann.com/gomorra"
	"math/rand"
	"log"
)

//TODO: make this channels or something cool
func randomFloat() float32 {
	return rand.Float32()
}

func main() {
	log.SetFlags(log.Lshortfile)

	localhost := gom.Remote{
		Hostname: "127.0.0.1:22",
	}

	err := localhost.Open("alex" , "/home/alex/.ssh/id_rsa")
	gom.FatalErr(err)

	fmt.Print(localhost.LsDir("/home/alex"))
}

func _main() {
	fmt.Println("Hello world")
	err := tui.Init()
	if err != nil {
		panic(err)
	}
	// don't forget to close
	defer tui.Close()

	theHello := tui.NewPar("Hello world in the TUI")
	theHello.Height = 3
	theHello.Width = 20
	theHello.TextFgColor = tui.ColorWhite
	theHello.BorderLabel = "Hola!"


	tui.Body.AddRows(
		tui.NewRow(
			tui.NewCol(12, 0, theHello),
			),
	)
	tui.Body.Align()

	tui.Render(tui.Body)


	tui.Handle("/sys/kbd/q", func(tui.Event) {
		// press q to quit
		tui.StopLoop()
	})

	tui.Loop()
}