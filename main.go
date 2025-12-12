package main

import (
	"fmt"
	"os"

	"github.com/HWCronicus/ssh-resume/src/models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(models.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
