package shortcuts

import (
	"fmt"
	"strings"
)

type AliasShortcut struct {
	Name    string
	Command string
}

func (as *AliasShortcut) GetName() string {
	return as.Name
}

func (as *AliasShortcut) GetCommand() string {
	return as.Command
}

func (as *AliasShortcut) GetText() string {
	return fmt.Sprintf("alias %v=%v\n", as.GetName(), as.GetCommand())
}

func (as *AliasShortcut) ParseText(line string) {
	trimmed := strings.TrimPrefix(line, "alias ")
	parts := strings.SplitN(trimmed, "=", 2)
	as.Name = parts[0]
	as.Command = strings.Trim(parts[1], `"`)
}
