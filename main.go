package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	mainView view = iota
	aboutMeView
	skillsView
	workExperienceView
	projectsView
	contactInfoView
)

type model struct {
	currentView   view
	width         int
	height        int
	cursor        int
	terminalWidth int
}

var (
	titleStyle = lipgloss.NewStyle().
			Border(lipgloss.Border{Bottom: "─"}, true).
			BorderForeground(lipgloss.Color("208")).
			Bold(true).
			Foreground(lipgloss.Color("208"))

	aboutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("178"))

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("208")).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

func initialModel() model {
	return model{
		currentView: mainView,
		cursor:      0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		if msg.Width > 200 {
			m.width = 200
		} else {
			m.width = msg.Width
		}
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
				if m.currentView != mainView {
					m.currentView = view(m.cursor + 1)
				}
			}

		case "right", "l":
			if m.cursor < 4 {
				m.cursor++
				if m.currentView != mainView {
					m.currentView = view(m.cursor + 1)
				}
			}

		case "enter":
			if m.currentView == mainView {
				m.currentView = view(m.cursor + 1)
			}

		case "b", "backspace", "esc":
			if m.currentView != mainView {
				m.currentView = mainView
			}

		case "1", "2", "3", "4", "5":
			selection := int(msg.String()[0] - '1')
			if selection >= 0 && selection < 5 {
				m.cursor = selection
				if m.currentView == mainView {
					m.currentView = view(selection + 1)
				} else {
					m.currentView = view(selection + 1)
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var content string

	switch m.currentView {
	case mainView:
		content = m.renderMainView("")
	case aboutMeView:
		content = m.renderDetailView("About Me", "This is the About Me view.")
	case skillsView:
		content = m.renderDetailView("Skills", "This is the Skills view.")
	case workExperienceView:
		content = m.renderDetailView("Work Experience", "This is the Work Experience view.")
	case projectsView:
		content = m.renderDetailView("Projects", "This is the Projects view.")
	case contactInfoView:
		content = m.renderDetailView("Contact Info", "This is the Contact Info view.")
	}

	// Apply border with adjusted size
	bordered := renderGradientBorder(content, m.width, m.height)

	// Center the content if terminal is wider than the view
	if m.terminalWidth > m.width {
		return lipgloss.NewStyle().
			Width(m.terminalWidth).
			AlignHorizontal(lipgloss.Center).
			Render(bordered)
	}

	return bordered
}

func (m model) RenderMainTitle() string {
	title :=
		"  ____ ____  _   _   ____\n" +
			" / ___/ ___|| | | | |  _ \\ ___  ___ _   _ _ __ ___   ___ \n" +
			" \\___ \\___ \\| |_| | | |_) / _ \\/ __| | | | '_ ` _ \\ / _ \\\n" +
			"  ___) |__) |  _  | |  _ <  __/\\__ \\ |_| | | | | | |  __/\n" +
			" |____/____/|_| |_| |_| \\_\\___||___/\\__,_|_| |_| |_|\\___|"

	return titleStyle.Render(title)
}

func (m model) RenderListItems() string {
	listItems := []string{
		"About Me",
		"Skills",
		"Work Experience",
		"Projects",
		"Contact Info",
	}
	var list string
	for i, item := range listItems {
		frontCursor := ">"
		backCursor := "<"
		if m.cursor == i {
			list += fmt.Sprintf("  %s %s %s  ", frontCursor, selectedItemStyle.Render(item), backCursor)
		} else {
			list += fmt.Sprintf("    %s    ", itemStyle.Render(item))
		}
	}
	return list
}

func (m model) RenderHelp() string {
	return helpStyle.Render("←/→: navigate • 1-5: quick select • enter: select • q: quit")
}

func (m model) renderMainView(middleContent string) string {

	title := m.RenderMainTitle()
	about := aboutStyle.Render("Welcome to Alan George's interactive resume! Navigate through the sections to learn more about me.")
	list := m.RenderListItems()
	help := m.RenderHelp()

	topContent := fmt.Sprintf("%s\n\n%s\n\n%s", title, about, list)

	top := lipgloss.NewStyle().
		Width(m.width).
		AlignHorizontal(lipgloss.Center).
		Render(topContent)

	var middle string
	if middleContent == "" {
		availableHeight := m.height - 8
		topContentHeight := lipgloss.Height(topContent)
		helpHeight := lipgloss.Height(help)
		spacerHeight := availableHeight - topContentHeight - helpHeight + 4
		if spacerHeight < 0 {
			spacerHeight = 0
		}

		middle = lipgloss.NewStyle().
			Height(spacerHeight).
			Render("")
	} else {
		middle = lipgloss.NewStyle().
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			Render(middleContent)
	}

	bottom := lipgloss.NewStyle().
		Width(m.width - 8).
		AlignHorizontal(lipgloss.Center).
		Render(help)

	return lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
}

func (m model) renderDetailView(title, description string) string {
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(1, 2).
		MarginTop(2).
		MarginBottom(2).
		Width(m.width - 15).
		Align(lipgloss.Center) // Make it full width minus padding/border

	viewContent := fmt.Sprintf("%s\n\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")).Render(title),
		itemStyle.Render(description))

	middleContent := contentStyle.Render(viewContent)

	return m.renderMainView(middleContent)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func getGradientColor(step, totalSteps int) lipgloss.Color {
	// Orange RGB: ~255, 135, 0 (approximately color 208)
	// Grey RGB: ~128, 128, 128 (approximately color 240)

	ratio := float64(step) / float64(totalSteps)

	// Interpolate between orange (208) and grey using ANSI color codes
	// Simplified: use predefined colors that form a gradient
	colors := []string{"208", "214", "180", "144", "245", "242", "240"}
	index := int(ratio * float64(len(colors)-1))
	if index >= len(colors) {
		index = len(colors) - 1
	}
	return lipgloss.Color(colors[index])
}

func renderGradientBorder(content string, width, height int) string {
	// Guard against invalid dimensions
	if width < 6 || height < 6 {
		return content
	}

	innerWidth := width - 2 // Account for left and right border characters

	lines := lipgloss.NewStyle().
		Width(innerWidth).
		Height(height-2). // Account for top and bottom border
		Padding(1, 2).
		Render(content)

	contentLines := strings.Split(lines, "\n")
	totalLines := len(contentLines)

	var result strings.Builder

	// Top border
	topColor := getGradientColor(0, totalLines)
	result.WriteString(lipgloss.NewStyle().Foreground(topColor).Render("╔" + strings.Repeat("═", innerWidth) + "╗"))

	// Side borders with gradient
	for i, line := range contentLines {
		color := getGradientColor(i+1, totalLines)
		borderStyle := lipgloss.NewStyle().Foreground(color)
		// Ensure line is exactly innerWidth
		paddedLine := lipgloss.NewStyle().Width(innerWidth).Render(line)
		result.WriteString("\n" + borderStyle.Render("║") + paddedLine + borderStyle.Render("║"))
	}

	// Bottom border
	result.WriteString("\n")
	bottomColor := getGradientColor(totalLines, totalLines)
	result.WriteString(lipgloss.NewStyle().Foreground(bottomColor).Render("╚" + strings.Repeat("═", innerWidth) + "╝"))

	return result.String()
}
