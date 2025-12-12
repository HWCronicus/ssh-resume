package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func MakeGradientStyles(colorA, colorB string, steps int) (s []lipgloss.Color) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0; i < steps; i++ {
		c := cA.BlendLuv(cB, float64(i)/float64(steps))
		s = append(s, lipgloss.Color(ColorToHex(c)))
	}
	return
}

func ColorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", ColorFloatToHex(c.R), ColorFloatToHex(c.G), ColorFloatToHex(c.B))
}

func ColorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func RenderGradientBorder(colorA, colorB, content string, width, height int) string {
	if width < 6 || height < 6 {
		return content
	}

	innerWidth := width - 2

	lines := lipgloss.NewStyle().
		Width(innerWidth).
		Height(height-2).
		Padding(1, 2).
		Render(content)

	contentLines := strings.Split(lines, "\n")
	totalLines := len(contentLines)

	var result strings.Builder

	topColor := MakeGradientStyles(colorA, colorB, totalLines)[0]
	result.WriteString(lipgloss.NewStyle().Foreground(topColor).Render("╔" + strings.Repeat("═", innerWidth) + "╗"))

	for i, line := range contentLines {
		color := MakeGradientStyles(colorA, colorB, totalLines)[i]
		borderStyle := lipgloss.NewStyle().Foreground(color)
		paddedLine := lipgloss.NewStyle().Width(innerWidth).Render(line)
		result.WriteString("\n" + borderStyle.Render("║") + paddedLine + borderStyle.Render("║"))
	}

	result.WriteString("\n")
	bottomColor := MakeGradientStyles(colorA, colorB, totalLines)[totalLines-1]
	result.WriteString(lipgloss.NewStyle().Foreground(bottomColor).Render("╚" + strings.Repeat("═", innerWidth) + "╝"))

	return result.String()
}
