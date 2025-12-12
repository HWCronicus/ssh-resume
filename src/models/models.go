package models

import (
	"fmt"

	"github.com/HWCronicus/ssh-resume/src/utils"
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

type Model struct {
	CurrentView   view
	Width         int
	Height        int
	Cursor        int
	TerminalWidth int
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

func InitialModel(height, width int) Model {
	return Model{
		CurrentView: mainView,
		Cursor:      0,
		Height:      height,
		Width:       width,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		if msg.Width > 200 {
			m.Width = 200
		} else {
			m.Width = msg.Width
		}
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.Cursor > 0 {
				m.Cursor--
				if m.CurrentView != mainView {
					m.CurrentView = view(m.Cursor + 1)
				}
			}

		case "right", "l":
			if m.Cursor < 4 {
				m.Cursor++
				if m.CurrentView != mainView {
					m.CurrentView = view(m.Cursor + 1)
				}
			}

		case "enter":
			if m.CurrentView == mainView {
				m.CurrentView = view(m.Cursor + 1)
			}

		case "b", "backspace", "esc":
			if m.CurrentView != mainView {
				m.CurrentView = mainView
			}

		case "1", "2", "3", "4", "5":
			selection := int(msg.String()[0] - '1')
			if selection >= 0 && selection < 5 {
				m.Cursor = selection
				if m.CurrentView == mainView {
					m.CurrentView = view(selection + 1)
				} else {
					m.CurrentView = view(selection + 1)
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var content string

	switch m.CurrentView {
	case mainView:
		content = m.renderMainView("")
	case aboutMeView:
		content = m.RenderDetailView("About Me", "This is the About Me view.")
	case skillsView:
		content = m.RenderDetailView("Skills", "This is the Skills view.")
	case workExperienceView:
		content = m.RenderDetailView("Work Experience", "This is the Work Experience view.")
	case projectsView:
		content = m.RenderDetailView("Projects", "This is the Projects view.")
	case contactInfoView:
		content = m.RenderDetailView("Contact Info", "This is the Contact Info view.")
	}

	bordered := utils.RenderGradientBorder("#ff5e00ff", "#555555ff", content, m.Width, m.Height)

	if m.TerminalWidth > m.Width {
		return lipgloss.NewStyle().
			Width(m.TerminalWidth).
			AlignHorizontal(lipgloss.Center).
			Render(bordered)
	}

	return bordered
}

func (m Model) RenderMainTitle() string {
	title :=
		"  ____ ____  _   _   ____\n" +
			" / ___/ ___|| | | | |  _ \\ ___  ___ _   _ _ __ ___   ___ \n" +
			" \\___ \\___ \\| |_| | | |_) / _ \\/ __| | | | '_ ` _ \\ / _ \\\n" +
			"  ___) |__) |  _  | |  _ <  __/\\__ \\ |_| | | | | | |  __/\n" +
			" |____/____/|_| |_| |_| \\_\\___||___/\\__,_|_| |_| |_|\\___|"

	return titleStyle.Render(title)
}

func (m Model) RenderListItems() string {
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
		if m.Cursor == i {
			list += fmt.Sprintf("  %s %s %s  ", frontCursor, selectedItemStyle.Render(item), backCursor)
		} else {
			list += fmt.Sprintf("    %s    ", itemStyle.Render(item))
		}
	}
	return list
}

func (m Model) RenderHelp() string {
	return helpStyle.Render("←/→: navigate • 1-5: quick select • enter: select • q: quit")
}

func (m Model) renderMainView(middleContent string) string {

	title := m.RenderMainTitle()
	about := aboutStyle.Render("Welcome to Alan George's interactive resume! Navigate through the sections to learn more about me.")
	list := m.RenderListItems()
	help := m.RenderHelp()

	topContent := fmt.Sprintf("%s\n\n%s\n\n%s", title, about, list)

	top := lipgloss.NewStyle().
		Width(m.Width).
		AlignHorizontal(lipgloss.Center).
		Render(topContent)

	var middle string
	if middleContent == "" {
		availableHeight := m.Height - 8
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
			Width(m.Width).
			AlignHorizontal(lipgloss.Center).
			Render(middleContent)
	}

	bottom := lipgloss.NewStyle().
		Width(m.Width - 8).
		AlignHorizontal(lipgloss.Center).
		Render(help)

	return lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
}

func (m Model) RenderDetailView(title, description string) string {
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(1, 2).
		MarginTop(2).
		MarginBottom(2).
		Width(m.Width - 15).
		Align(lipgloss.Center) // Make it full width minus padding/border

	viewContent := fmt.Sprintf("%s\n\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")).Render(title),
		itemStyle.Render(description))

	middleContent := contentStyle.Render(viewContent)

	return m.renderMainView(middleContent)
}
