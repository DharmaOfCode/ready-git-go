package main

import (
	"fmt"
	"os"
	"flag"
	"os/signal"
	"sync"
	"io/ioutil"
	"log"
	"os/exec"
)

type State struct {
	Path         string
	Username     string
	Password     string
	NewOrigin    string
	ChangeOrigin bool
	Pull         bool
	Verbose      bool
	StdIn        bool
	Quiet        bool
	Threads      int
	Terminate	 bool
	SignalChan   chan os.Signal
}

func ParseCmdLine() *State {
	valid := true

	s := State{
		StdIn: false,
	}
	s.Threads = 10
	flag.StringVar(&s.Path, "p", "", "Repositories file path")
	flag.StringVar(&s.Username, "U", "", "Username for git Auth")
	flag.StringVar(&s.Password, "P", "", "Password for git Auth")
	flag.BoolVar(&s.Pull, "u", false, "Update all repositories in file path (must provide file path)")
	flag.BoolVar(&s.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&s.Quiet, "q", false, "Quit output")
	flag.Parse()

	Banner(&s)

	if s.Path == "" {
		fmt.Println("[!] File Path (-p): Must be specified")
		valid = false
	}

	if s.Path[len(s.Path)-1:] != "\\" || s.Path[len(s.Path)-1:] == "/"{
		s.Path += "\\"
	}

	if valid {
		return &s
	} else {
		Ruler(&s)
	}

	return nil
}

func Ruler(s *State) {
	if !s.Quiet {
		fmt.Println("==============================================================")
	}
}

func Banner(state *State) {
	if state.Quiet {
		return
	}

	fmt.Println("")
	fmt.Println("readygitgo			By Alex Useche")
	Ruler(state)
}

func PrepareSignalHandler(s *State) {
	s.SignalChan = make(chan os.Signal, 1)
	signal.Notify(s.SignalChan, os.Interrupt)
	go func() {
		for _ = range s.SignalChan {
			// caught CTRL+C
			if !s.Quiet {
				fmt.Println("[!] Keyboard interrup detected. Terminating...")
				s.Terminate = true
				panic("[!] Panic! Panic! Panic!")
			}
		}
	}()

}

func Process(s *State) {
	PrepareSignalHandler(s)

	dirChan := make(chan string, s.Threads)

	processorGroup := new(sync.WaitGroup)
	processorGroup.Add(s.Threads)
	printerGroup := new(sync.WaitGroup)
	printerGroup.Add(1)

	// Create goroutines for each thread
	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				dir := <-dirChan

				// Are all directories traversed?
				if dir == "" {
					break
				}

				// Mode-specific processing
				//s.Processor(s, dir)
			}

			// Indicate to the wait group tha the thread has finished
			processorGroup.Done()
		}()
	}

	folders, err := ioutil.ReadDir(s.Path)
	if err != nil {
		log.Fatal(err)
	}

	//pwd, _ := os.Getwd()
	for _, f := range folders {
		if f.IsDir(){
			gitpath := s.Path + f.Name()
			fmt.Println(s.Path + f.Name())
			out, err := exec.Command("git", "-C", gitpath, "log").
				Output()
			if err != nil{
				fmt.Println(out)
			}
		}
	}
}

func main() {
	state := ParseCmdLine()
	if state != nil {
		Process(state)
	}
}
