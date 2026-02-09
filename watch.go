// Copyright (c) 2026 Eric Fialkowski
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/buger/goterm"
	"github.com/mitchellh/go-ps"
	"github.com/tevino/abool"
)

var interval = flag.Int("n", 5, "Interval in seconds")
var runWithCommand = flag.Bool("x", false, "Run with command processor")
var hideTitle = flag.Bool("t", false, "Hide title bar")
var exitOnError = flag.Bool("e", false, "Exit on non-zero return of command")
var preciseInterval = flag.Bool("p", false, "Try to run at precise interval")
var exitOnChange = flag.Bool("g", false, "Exit when output changes")
var oldOutput string

var once sync.Once
var cond = abool.New()

func init() {
	flag.Parse()
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("Use the native watch command instead of this")
		os.Exit(5)
	}
	cmdArray := flag.Args()
	if len(cmdArray) == 0 {
		fmt.Println("Must include a command to run")
		flag.Usage()
		os.Exit(3)
	}

	var cmd string
	var cmdArgs []string

	if *runWithCommand {
		c, e := ps.FindProcess(os.Getppid())
		if e != nil {
			fmt.Printf("Error getting parent command processor: %v\n", e)
			os.Exit(4)
		}
		cmd = c.Executable()
		cmdArgs = make([]string, len(cmdArgs)+1)
		if cmd == "cmd.exe" {
			cmdArgs[0] = "/c"
		} else {
			cmdArgs[0] = "-c"
		}
		cmdArgs = append(cmdArgs, cmdArray...)
	} else {
		cmd = cmdArray[0]
		cmdArgs = cmdArray[1:]
	}

	run(time.Now(), cmd, cmdArgs)
	nextRun := time.Now().Add(time.Duration(*interval) * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			if time.Now().After(nextRun) || time.Now().Equal(nextRun) {
				startOfRun := time.Now()
				run(t, cmd, cmdArgs)
				if *preciseInterval {
					nextRun = startOfRun.Add(time.Duration(*interval) * time.Second)
				} else {
					nextRun = time.Now().Add(time.Duration(*interval) * time.Second)
				}
			}
		}
	}()

	select {}
}

func run(t time.Time, name string, args []string) {
	if cond.SetToIf(false, true) {
		goterm.Clear()
		goterm.MoveCursor(1, 1)
		if !*hideTitle {
			goterm.Printf("Every %ds: %s %s", *interval, name, strings.Join(args, " "))
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

			goterm.MoveCursor(1, 3)
		}
		cmd := exec.Command(name, args...)
		output, err := cmd.CombinedOutput()
		if err == nil {
			goterm.Print(string(output))
		} else {
			goterm.Printf("Error: %v", err)
			if *exitOnError {
				goterm.Flush()
				os.Exit(1)
			}
		}
		goterm.Flush()
		if *exitOnChange {
			once.Do(func() {
				oldOutput = string(output)
			})
			if oldOutput != string(output) {
				os.Exit(2)
			}
		}

		cond.UnSet()
	}
}
