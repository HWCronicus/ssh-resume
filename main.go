package main

import (
	"fmt"
	"os"

	"github.com/HWCronicus/ssh-resume/src/models"
	"github.com/HWCronicus/ssh-resume/src/server"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

func main() {
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width, height = 200, 50
	}
	p := tea.NewProgram(models.InitialModel(height, width), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	server.StartServer()
}
