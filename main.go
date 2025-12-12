package main

import (
	"github.com/HWCronicus/ssh-resume/src/server"
)

func main() {
	// width, height, err := term.GetSize(os.Stdout.Fd())
	// if err != nil {
	// 	width, height = 200, 300
	// }
	// p := tea.NewProgram(models.InitialModel(height, width), tea.WithAltScreen())

	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	os.Exit(1)
	// }
	server.StartServer()
}
