package main

import (
	"flag"
	"fmt"
	"os/exec"
	"rcsc/pkg/shortcuts"
	"strings"
)

func main() {
	deleteFlag := flag.Bool("d", false, "Delete specified Shortcut")
	helpFlag := flag.Bool("help", false, "Display help")
	helpFlagShort := flag.Bool("h", false, "Display help")
	initFlag := flag.Bool("init", false, "Init rcsc to work on current shell")
	flag.Parse()

	if *helpFlag || *helpFlagShort {
		printHelp()
		return
	}

	args := flag.Args()
	manager := shortcuts.ShortcutsManager{}

	if *initFlag {
		manager.InitRc()
		return
	} else {
		manager.ParseRc()
	}
	changeMade := false

	if *deleteFlag {
		if len(args) != 1 {
			fmt.Println("To delete a shortcut:")
			fmt.Println("$ rcsc -d <name_of_shortcut>")
			printHelp()
			return
		}
		manager.DeleteShortcut(args[0])
		changeMade = true
	} else if len(args) == 0 {
		manager.ListShortcuts()
	} else {
		manager.AddShortcut(args[0], strings.Join(args[1:], " "))
		changeMade = true
	}
	if changeMade {
		manager.UpdateRc()
		// Run the source command
		cmd := exec.Command("zsh", "-c", "source ~/.zshrc")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error sourcing .zshrc: %v\n", err)
			return
		}
	}
}

func printHelp() {
	fmt.Println("\nUsage:")
	fmt.Println("$ rcsc                                                        List all shortcuts")
	fmt.Println("$ rcsc <name_of_shortcut>                                     Run the specified shortcut")
	fmt.Println("$ rcsc <name_of_shortcut> full command                        Add a shortcut")
	fmt.Println("$ rcsc <name_of_shortcut> $param1 - full $param1 command      Add a shortcut with params (Seperate with -)")
	fmt.Println("$ rcsc -d <name_of_shortcut>                                  Delete the specified shortcut")
	fmt.Println("$ rcsc --init                                                 Display this help message")
	fmt.Println("$ rcsc --help, -h                                             Init rcsc to work on current shell")
}
