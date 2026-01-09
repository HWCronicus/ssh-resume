package models

import (
	"fmt"
	"strings"

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
				Background(lipgloss.Color("208")).
				Foreground(lipgloss.Color("0")).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	highlightColor    = lipgloss.Color("208")
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 2).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

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
		m.Width = min(msg.Width, 200)
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "right", "l":
			if m.Cursor < 4 {
				m.Cursor++
			}

		case "enter":
			m.CurrentView = view(m.Cursor + 1)

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

	bordered := utils.RenderGradientBorder("#ff7300", "#666666", content, m.Width, m.Height)

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

func (m Model) RenderTabs(content string) string {
	tabs := []string{"About Me", "Skills", "Work Experience", "Projects", "Contact Info"}

	var renderedTabs []string
	for i, t := range tabs {
		var style lipgloss.Style
		isActive := i == m.Cursor
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	contentWidth := min(200, m.Width-20)
	totalWidth := contentWidth + windowStyle.GetHorizontalFrameSize() - 4

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	tabsWidth := lipgloss.Width(tabRow)

	leftGap := (totalWidth - tabsWidth) / 2
	rightGap := totalWidth - tabsWidth - leftGap

	var leftBorder string
	if leftGap > 0 {
		leftBorder = lipgloss.NewStyle().
			Foreground(highlightColor).
			Render("┌" + strings.Repeat("─", leftGap-1))
	}

	var rightBorder string
	if rightGap > 0 {
		rightBorder = lipgloss.NewStyle().
			Foreground(highlightColor).
			Render(strings.Repeat("─", rightGap-1) + "┐")
	}

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, leftBorder, tabRow, rightBorder)

	doc := strings.Builder{}
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(contentWidth).Render(content))

	return doc.String()
}

func (m Model) RenderHelp() string {
	return helpStyle.Render("←/→: navigate • 1-5: quick select • enter: select • q: quit")
}

func (m Model) renderMainView(middleContent string) string {

	title := m.RenderMainTitle()
	about := aboutStyle.Render("Welcome to Alan George's interactive resume! Navigate through the sections to learn more about me.")

	var tabs string
	if middleContent == "" {
		tabs = m.RenderTabs("Select a section to view details")
	} else {
		tabs = m.RenderTabs(middleContent)
	}

	help := m.RenderHelp()

	topContent := fmt.Sprintf("%s\n\n%s\n\n", title, about)

	top := lipgloss.NewStyle().
		Width(m.Width - 8).
		AlignHorizontal(lipgloss.Center).
		Render(topContent)

	availableHeight := m.Height - 8
	topContentHeight := lipgloss.Height(topContent)
	helpHeight := lipgloss.Height(help)
	spacerHeight := availableHeight - topContentHeight - helpHeight + 4
	if spacerHeight < 0 {
		spacerHeight = 0
	}

	middle := lipgloss.NewStyle().
		Height(spacerHeight).
		Width(m.Width - 8).
		PaddingLeft(3).
		AlignHorizontal(lipgloss.Center).
		Render(tabs)

	bottom := lipgloss.NewStyle().
		Width(m.Width - 8).
		AlignHorizontal(lipgloss.Center).
		Render(help)

	return lipgloss.JoinVertical(lipgloss.Top, top, middle, bottom)
}

func (m Model) RenderDetailView(titleText, descriptionText string) string {
	contentStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(min(200, m.Width-15))

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")).Render(titleText)
	description := itemStyle.Render(descriptionText)

	viewContent := fmt.Sprintf("%s\n\n%s",
		title,
		description)

	middleContent := contentStyle.Render(viewContent)

	return m.renderMainView(middleContent)
}
