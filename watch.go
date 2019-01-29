package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/buger/goterm"
)

var interval = flag.Int("n", 5, "Interval in seconds")

func init() {
	flag.Parse()
}

func main() {
	cmdArray := flag.Args()
	if len(cmdArray) == 0 {
		fmt.Println("Must include a command to run")
		os.Exit(1)
	}
	cmd := cmdArray[0]
	cmdArgs := cmdArray[1:]

	run(time.Now(), cmd, cmdArgs)
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	go func() {
		for t := range ticker.C {
			run(t, cmd, cmdArgs)
		}
	}()

	select {}
}

func run(t time.Time, name string, cmdArgs []string) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Printf("Every %ds: %s %s", *interval, name, strings.Join(cmdArgs, " "))
	width := goterm.Width()
	ts := t.Format("Mon Jan _2 15:04:05 2006")
	hn, err := os.Hostname()
	if err == nil {
		s := fmt.Sprintf("%s: %s", hn, ts)
		goterm.MoveCursor(width-len(s), 1)
		goterm.Print(s)
	} else {
		goterm.MoveCursor(width-len(ts), 1)
		goterm.Print(ts)
	}

	goterm.MoveCursor(2, 3)
	cmd := exec.Command(name, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err == nil {
		goterm.Print(string(output))
	} else {
		goterm.Printf("Error: %v", err)
	}
	goterm.Flush()
}
