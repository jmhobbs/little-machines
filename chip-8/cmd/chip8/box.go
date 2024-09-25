package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// From https://gist.github.com/meowgorithm/1777377a43373f563476a2bcb7d89306

type BoxWithLabel struct {
	BoxStyle   lipgloss.Style
	LabelStyle lipgloss.Style
}

func NewDefaultBoxWithLabel() BoxWithLabel {
	return BoxWithLabel{
		BoxStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			PaddingLeft(1),

		LabelStyle: lipgloss.NewStyle().
			PaddingTop(0).
			PaddingBottom(0).
			PaddingLeft(1).
			PaddingRight(1),
	}
}

func (b BoxWithLabel) Render(label, content string, width int) string {
	var (
		// Query the box style for some of its border properties so we can
		// essentially take the top border apart and put it around the label.
		border          lipgloss.Border             = b.BoxStyle.GetBorderStyle()
		topBorderStyler func(strs ...string) string = lipgloss.NewStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		topLeft         string                      = topBorderStyler(border.TopLeft)
		topRight        string                      = topBorderStyler(border.TopRight)

		renderedLabel string = b.LabelStyle.Render(label)
	)

	// Render top row with the label
	borderWidth := b.BoxStyle.GetHorizontalBorderSize()
	cellsShort := max(0, width+borderWidth-lipgloss.Width(topLeft+topRight+renderedLabel))
	gap := strings.Repeat(border.Top, cellsShort)
	top := topLeft + renderedLabel + topBorderStyler(gap) + topRight

	// Render the rest of the box
	bottom := b.BoxStyle.Copy().
		BorderTop(false).
		Width(width).
		Render(content)

	// Stack the pieces
	return top + "\n" + bottom
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
