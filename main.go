package main

import (
	"fmt"
	"os"
	"flag"
	"io/ioutil"
	"log"
	"os/exec"
)

type State struct {
	Path         	string
	Username     	string
	Password     	string
	NewOriginBase   string
	ChangeOrigin 	bool
	Pull         	bool
	Verbose      	bool
	StdIn        	bool
	Quiet        	bool
	Threads      	int
	Terminate	 	bool
	SignalChan   	chan os.Signal
}

func ParseCmdLine() *State {
	valid := true

	s := State{
		StdIn: false,
	}
	s.Threads = 10
	flag.StringVar(&s.Path, "p", "", "Repositories file path")
	flag.StringVar(&s.NewOriginBase, "o", "", "New origin base file path")
	flag.StringVar(&s.Username, "U", "", "Username for git Auth")
	flag.StringVar(&s.Password, "P", "", "Password for git Auth")
	flag.BoolVar(&s.Pull, "u", false, "Update all repositories in file path (must provide file path)")
	flag.BoolVar(&s.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&s.Quiet, "q", false, "Quit output")
	flag.BoolVar(&s.ChangeOrigin, "c", false, "Change git repo origin")
	flag.Parse()

	Banner(&s)

	if s.Path == "" {
		fmt.Println("[!] File Path (-p): Must be specified")
		valid = false
	}

	if s.Path[len(s.Path)-1:] != "\\" || s.Path[len(s.Path)-1:] == "/"{
		s.Path += "\\"
	}

	if s.ChangeOrigin{
		if s.NewOriginBase == ""{
			fmt.Println("[!] You must provide the new base URL for the origin")
			valid = false
		} else {
			if s.NewOriginBase[len(s.NewOriginBase)-1:] != "/"{
				s.NewOriginBase += "/"
			}
		}
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

func Process(s *State) {

	if s.ChangeOrigin{
		UpdateOrigins(s)
	}

	if s.Pull{
		PullAllRepos(s)
	}
}

func PullAllRepos(s *State){
	folders, err := ioutil.ReadDir(s.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range folders {
		if f.IsDir(){
			gitpath := s.Path + f.Name()

			if s.Verbose{
				fmt.Println("Updating repo " + f.Name())
				fmt.Println("running command: " + "git -C " + "pull" + gitpath )
			}
			out, err := exec.Command("git", "-C", gitpath, "pull").
				Output()
			if err != nil{
				log.Fatal(err)
			}
			fmt.Printf("%s\n", out)
			Ruler(s)
		}
	}
}

func UpdateOrigins(s *State){

	folders, err := ioutil.ReadDir(s.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range folders {
		if f.IsDir(){
			gitpath := s.Path + f.Name()

			origin := s.NewOriginBase + f.Name()
			if s.Verbose{
				fmt.Println("Updating remote origin for " + f.Name())
				fmt.Println("running command: " + "git -C " + gitpath + " remote set-url orign")
			}
			out, err := exec.Command("git", "-C", gitpath, "remote", "set-url", "origin", origin).
				Output()
			if err != nil{
				log.Fatal(err)
			}
			fmt.Printf("%s\n", out)
			Ruler(s)
		}
	}

	Ruler(s)
}

func main() {
	state := ParseCmdLine()
	if state != nil {
		Process(state)
	}
}
