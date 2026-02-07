package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	title        lipgloss.Style
	meta         lipgloss.Style
	header       lipgloss.Style
	headerText   lipgloss.Style
	headerMeta   lipgloss.Style
	summaryBar   lipgloss.Style
	summaryText  lipgloss.Style
	checkList    lipgloss.Style
	footer       lipgloss.Style
	footerText   lipgloss.Style
	empty        lipgloss.Style
	scrollInfo   lipgloss.Style
	columnHeader lipgloss.Style
	divider      lipgloss.Style
	colorPass    lipgloss.Style
	colorFail    lipgloss.Style
	colorWarn    lipgloss.Style
	colorUnknown lipgloss.Style
	colorRunning lipgloss.Style
	colorQueued  lipgloss.Style
	colorIdle    lipgloss.Style
	width        int
}

func newStyles(width int) styles {
	headerBg := lipgloss.Color("#235")
	summaryBg := lipgloss.Color("#61AFEF")
	footerBg := lipgloss.Color("#235")

	return styles{
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")),

		meta: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")),

		header: lipgloss.NewStyle().
			Background(headerBg).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1).
			Width(width),

		headerText: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			Background(headerBg),

		headerMeta: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Background(headerBg),

		summaryBar: lipgloss.NewStyle().
			Background(summaryBg).
			Foreground(lipgloss.Color("0")).
			Padding(0, 1).
			Width(width),

		summaryText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(summaryBg),

		checkList: lipgloss.NewStyle().
			Padding(0, 1).
			Width(width),

		columnHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("252")),

		divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color("238")),

		footer: lipgloss.NewStyle().
			Background(footerBg).
			Foreground(lipgloss.Color("243")).
			Padding(0, 1).
			Width(width),

		footerText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Background(footerBg),

		empty: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true),

		scrollInfo: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Align(lipgloss.Right),

		colorPass: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")),

		colorFail: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E06C75")),

		colorWarn: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5C07B")),

		colorUnknown: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),

		colorRunning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61AFEF")),

		colorQueued: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C678DD")),

		colorIdle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),

		width: width,
	}
}
