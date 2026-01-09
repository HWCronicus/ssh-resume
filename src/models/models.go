package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/HWCronicus/ssh-resume/src/utils"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
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

const (
	minWidth  = 150
	minHeight = 50
)

type Model struct {
	CurrentView   view
	Width         int
	Height        int
	Cursor        int
	TerminalWidth int
	Viewport      viewport.Model
	Ready         bool
}

var (
	titleStyle = lipgloss.NewStyle().
			Border(lipgloss.Border{Bottom: "─"}, true).
			BorderForeground(lipgloss.Color("208")).
			Foreground(lipgloss.Color("208"))

	aboutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	highlightColor    = lipgloss.Color("208")
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true).Foreground(lipgloss.Color("208")).Underline(true).Bold(true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 2).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop().UnsetBorderBottom()

	infoStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 1)
	}()
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
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TerminalWidth = msg.Width
		m.Width = min(msg.Width, 150)
		m.Height = min(msg.Height, 50)

		if !m.Ready {
			contentWidth := min(150, m.Width-20)
			viewportWidth := contentWidth - 4
			viewportHeight := m.Height - 30
			m.Viewport = viewport.New(viewportWidth, viewportHeight)
			m.Viewport.YPosition = 0
			m.Ready = true
		} else {
			contentWidth := min(150, m.Width-20)
			viewportWidth := contentWidth - 4
			viewportHeight := m.Height - 30
			m.Viewport.Width = viewportWidth
			m.Viewport.Height = viewportHeight
		}
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
			m.loadViewportContent()

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
				m.loadViewportContent()
			}
		}
	}

	if m.CurrentView != mainView && m.Ready {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.TerminalWidth < minWidth || m.Height < minHeight {
		return m.renderResizeMessage()
	}

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

func (m Model) renderResizeMessage() string {
	messageStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("208")).
		Align(lipgloss.Center).
		Width(m.TerminalWidth).
		Height(m.Height)

	message := fmt.Sprintf(
		"Your Terminal window is too small.\n\n"+
			"For an optimal experience, please resize your terminal window.\n\n"+
			"Current size: %dx%d\n"+
			"Minimum size: %dx%d\n\n",
		m.TerminalWidth, m.Height,
		minWidth, minHeight,
	)

	return messageStyle.Render(message)
}

func (m Model) RenderMainTitle() string {
	title :=

		"  █████████   ████                         █████████                                                   ██████████                        \n" +
			"  ███░░░░░███ ░░███                        ███░░░░░███                                                 ░░███░░░░███                      \n" +
			" ░███    ░███  ░███   ██████   ████████   ███     ░░░   ██████   ██████  ████████   ███████  ██████     ░███   ░░███  ██████  █████ █████\n" +
			" ░███████████  ░███  ░░░░░███ ░░███░░███ ░███          ███░░███ ███░░███░░███░░███ ███░░███ ███░░███    ░███    ░███ ███░░███░░███ ░░███ \n" +
			" ░███░░░░░███  ░███   ███████  ░███ ░███ ░███    █████░███████ ░███ ░███ ░███ ░░░ ░███ ░███░███████     ░███    ░███░███████  ░███  ░███ \n" +
			" ░███    ░███  ░███  ███░░███  ░███ ░███ ░░███  ░░███ ░███░░░  ░███ ░███ ░███     ░███ ░███░███░░░      ░███    ███ ░███░░░   ░░███ ███  \n" +
			" █████   █████ █████░░████████ ████ █████ ░░█████████ ░░██████ ░░██████  █████    ░░███████░░██████  ██ ██████████  ░░██████   ░░█████   \n" +
			"░░░░░   ░░░░░ ░░░░░  ░░░░░░░░ ░░░░ ░░░░░   ░░░░░░░░░   ░░░░░░   ░░░░░░  ░░░░░      ░░░░░███ ░░░░░░  ░░ ░░░░░░░░░░    ░░░░░░     ░░░░░    \n" +
			"                                                                                   ███ ░███                                              \n" +
			"                                                                                  ░░██████                                               \n" +
			"                                                                                   ░░░░░░                                                \n"

	return titleStyle.Render(title)
}

func (m Model) RenderTabs(content string, showFooter bool) string {
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

	contentWidth := min(150, m.Width-20)
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

	if showFooter {
		footer := m.footerView()
		if footer != "" {
			doc.WriteString("\n")
			doc.WriteString(footer)
		}
	} else {
		bottomBorder := lipgloss.NewStyle().
			Foreground(highlightColor).
			Render("└" + strings.Repeat("─", contentWidth) + "┘")
		doc.WriteString("\n")
		doc.WriteString(bottomBorder)
	}

	return doc.String()
}

func (m Model) RenderHelp() string {
	return helpStyle.Render("←/→: navigate • ↑/↓: scroll • 1-5: quick select • enter: select • q: quit")
}

func (m Model) footerView() string {
	if !m.Ready || m.CurrentView == mainView {
		return ""
	}

	contentWidth := min(150, m.Width-20)

	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	infoWidth := lipgloss.Width(info)

	leftLineWidth := (contentWidth - infoWidth) / 2
	if leftLineWidth < 1 {
		leftLineWidth = 1
	}

	rightLineWidth := contentWidth - infoWidth - leftLineWidth
	if rightLineWidth < 1 {
		rightLineWidth = 1
	}

	leftLine := lipgloss.NewStyle().
		Foreground(highlightColor).
		Render("└" + strings.Repeat("─", leftLineWidth-1) + "┤")

	rightLine := lipgloss.NewStyle().
		Foreground(highlightColor).
		Render("├" + strings.Repeat("─", rightLineWidth-1) + "┘")

	return lipgloss.JoinHorizontal(lipgloss.Center, leftLine, info, rightLine)
}

func (m Model) renderMainView(middleContent string) string {

	title := m.RenderMainTitle()
	about := aboutStyle.Render("Welcome to Terminal based version of AlanGeorge.Dev, navigate through the sections to learn more about me.")

	var tabs string
	showFooter := m.CurrentView != mainView && middleContent != ""
	if middleContent == "" {
		tabs = m.RenderTabs("Select a section to view details", false)
	} else {
		tabs = m.RenderTabs(middleContent, showFooter)
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
		Height(spacerHeight - 10).
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

func (m *Model) loadViewportContent() tea.Cmd {
	if !m.Ready {
		return nil
	}

	var filename string
	switch m.CurrentView {
	case aboutMeView:
		filename = "content/about-me.md"
	case skillsView:
		filename = "content/skills.md"
	case workExperienceView:
		filename = "content/work-experience.md"
	case projectsView:
		filename = "content/projects.md"
	case contactInfoView:
		filename = "content/contact-info.md"
	default:
		return nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		errorMsg := fmt.Sprintf("File not found: %s\n\nError: %s\n\nPlease make sure the file exists.", filename, err.Error())
		m.Viewport.SetContent(errorMsg)
		return nil
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.Viewport.Width-2),
	)
	if err != nil {
		m.Viewport.SetContent(string(content))
		return nil
	}

	renderedContent, err := renderer.Render(string(content))
	if err != nil {
		m.Viewport.SetContent(string(content))
		return nil
	}

	m.Viewport.SetContent(renderedContent)
	m.Viewport.GotoTop()
	return nil
}

func (m Model) RenderDetailView(titleText, descriptionText string) string {
	if m.TerminalWidth < minWidth || m.Height < minHeight {
		return m.renderResizeMessage()
	}

	var viewContent string
	if m.Ready {
		viewContent = m.Viewport.View()
		if viewContent == "" {
			viewContent = "Press a number key (1-5) or use arrow keys and Enter to select a section."
		}
	} else {
		viewContent = "Initializing..."
	}

	contentStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Width(min(150, m.Width-15))

	middleContent := contentStyle.Render(viewContent)

	return m.renderMainView(middleContent)
}
