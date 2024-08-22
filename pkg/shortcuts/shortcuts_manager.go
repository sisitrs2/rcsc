package shortcuts

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type ShortcutsManager struct {
	Shortcuts []Shortcut
}

func (sm *ShortcutsManager) DeleteShortcut(name string) {
	for i, shortcut := range sm.Shortcuts {
		if shortcut.GetName() == name {
			sm.Shortcuts = append(sm.Shortcuts[:i], sm.Shortcuts[i+1:]...)
			fmt.Printf("Deleted shortcut %v\n", name)
			return
		}
	}
}

func (sm *ShortcutsManager) AddShortcut(name string, command string) {
	sc := &AliasShortcut{
		Name:    name,
		Command: command,
	}
	sm.Shortcuts = append(sm.Shortcuts, sc)
}

func (sm *ShortcutsManager) PrintShortcuts() {
	for _, sc := range sm.Shortcuts {
		fmt.Print(sc.GetText())
	}
}

func (sm *ShortcutsManager) GetShortcutsText() string {
	text := ""
	for _, sc := range sm.Shortcuts {
		text += sc.GetText()
	}
	return text
}

func (sm *ShortcutsManager) ListShortcuts() {
	// Find the length of the longest shortcut name
	maxNameLen := 0
	for _, sc := range sm.Shortcuts {
		if len(sc.GetName()) > maxNameLen {
			maxNameLen = len(sc.GetName())
		}
	}

	// Print each shortcut with aligned output
	for _, sc := range sm.Shortcuts {
		name := sc.GetName()
		command := sc.GetCommand()
		// Format the output so that the `:` aligns
		fmt.Printf("%-*s : %v\n", maxNameLen, name, command)
	}
}

func (sm *ShortcutsManager) ParseRc() error {
	shell := os.Getenv("SHELL")
	switch shell {
	case "/bin/zsh":
		sm.parseRcFile("~/.zshrc")
	case "/bin/bash":
		sm.parseRcFile("~/.bashrc")
	case "/bin/sh":
		// sm.parseRcFile("~/.profile")
		sm.parseRcFile("~/.zshrc")
	default:
		return fmt.Errorf("ERROR: Unrecognized shell %s", shell)
	}
	return nil
}

func (sm *ShortcutsManager) initRcFile(fname string) {
	path, _ := expandPath(fname)
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening rc file %v\n", fname)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isRCSC := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "#>>> RCSC Section <<<#" {
			isRCSC = true
		} else if isRCSC && line == "#>>> End RCSC <<<#" {
			fmt.Println("RCSC is already initialized.")
			return
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
	}
	if isRCSC {
		fmt.Printf("File %v is corrupt - start of RCSC section discoverd but not end.\nPlease remove it manually and init again.\n", fname)
		return
	}

	file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("error opening rc file for writing %v\n", fname)
		return
	}
	defer file.Close()

	_, _ = file.WriteString("#>>> RCSC Section <<<#" + "\n")
	_, _ = file.WriteString("#>>> End RCSC <<<#" + "\n")
	println("Succesfully initialized.")
}

func (sm *ShortcutsManager) InitRc() error {
	shell := os.Getenv("SHELL")
	switch shell {
	case "/bin/zsh":
		sm.initRcFile("~/.zshrc")
	case "/bin/bash":
		sm.initRcFile("~/.bashrc")
	case "/bin/sh":
		sm.initRcFile("~/.profile")
	default:
		return fmt.Errorf("ERROR: Unrecognized shell %s", shell)
	}
	return nil
}

func (sm *ShortcutsManager) UpdateRc() error {
	shell := os.Getenv("SHELL")
	switch shell {
	case "/bin/zsh":
		sm.updateRcFile("~/.zshrc")
	case "/bin/bash":
		sm.updateRcFile("~/.bashrc")
	case "/bin/sh":
		// sm.updateRcFile("~/.profile")
		sm.updateRcFile("~/.zshrc")
	default:
		return fmt.Errorf("ERROR: Unrecognized shell %s", shell)
	}
	return nil
}

func (sm *ShortcutsManager) updateRcFile(fname string) {
	path, _ := expandPath(fname)
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening rc file %v\n", fname)
	}
	defer file.Close()

	inRCSC := false
	var fileLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !inRCSC {
			fileLines = append(fileLines, line)
		}
		if line == "#>>> RCSC Section <<<#" {
			inRCSC = true
		} else if line == "#>>> End RCSC <<<#" {
			inRCSC = false
			fileLines = append(fileLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
	}

	file, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("error opening rc file for writing %v\n", fname)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	inRCSC = false
	for _, line := range fileLines {
		_, _ = writer.WriteString(line + "\n")
		if line == "#>>> RCSC Section <<<#" {
			// Add shortcuts
			_, _ = writer.WriteString(sm.GetShortcutsText())
		}
	}
	writer.Flush()

}

func (sm *ShortcutsManager) parseRcFile(fname string) {
	path, _ := expandPath(fname)
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening rc file %v\n", fname)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inRCSC := false
	cantFindRCSC := true
	line := ""
	scs := []Shortcut{}
	for scanner.Scan() {
		line = scanner.Text()
		if inRCSC && len(line) > 3 && line[0] != '#' {
			scs = append(scs, addShortcut(line))
		}
		if line == "#>>> RCSC Section <<<#" {
			inRCSC = true
			cantFindRCSC = false
		} else if line == "#>>> End RCSC <<<#" {
			inRCSC = false
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
	}

	if cantFindRCSC {
		fmt.Println("RCSC has never been initialized on the shell.")
		fmt.Println("Please run:\n$ rcsc --init")
	} else {
		sm.Shortcuts = scs
	}

}

func addShortcut(line string) Shortcut {
	words := strings.Fields(line)
	switch words[0] {
	case "alias":
		sc := &AliasShortcut{}
		sc.ParseText(line)
		return sc
	}
	return nil
}

func expandPath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, path[1:])
	}
	return path, nil
}
