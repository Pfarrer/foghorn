package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/pfarrer/foghorn/scheduler"
)

type model struct {
	scheduler    *scheduler.Scheduler
	logLevel     string
	uptime       time.Time
	width        int
	height       int
	maxCheckRows int
}

type tickMsg time.Time

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func NewModel(sched *scheduler.Scheduler, logLevel string) model {
	return model{
		scheduler:    sched,
		logLevel:     logLevel,
		uptime:       sched.GetStartTime(),
		maxCheckRows: 20,
	}
}

func (m model) Init() tea.Cmd {
	return tickEvery(time.Second)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		return m, tickEvery(time.Second)
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	styles := newStyles(m.width)

	var builder strings.Builder

	builder.WriteString(m.renderHeader(styles))
	builder.WriteString("\n")
	builder.WriteString(m.renderSummaryBar(styles))
	builder.WriteString("\n")
	builder.WriteString(m.renderCheckList(styles))
	builder.WriteString("\n")
	builder.WriteString(m.renderFooter(styles))

	return builder.String()
}

func (m model) getAvailableCheckRows() int {
	available := m.height - 6
	if available < 1 {
		return 1
	}
	return available
}

func (m model) renderHeader(styles styles) string {
	uptime := time.Since(m.uptime).Round(time.Second)
	title := styles.headerText.Render("Foghorn")
	uptimeStr := styles.headerMeta.Render(fmt.Sprintf("Uptime: %s", uptime))
	levelStr := styles.headerMeta.Render(fmt.Sprintf("Log Level: %s", m.logLevel))
	separator := styles.headerMeta.Render("  ")

	left := lipgloss.JoinHorizontal(lipgloss.Top, title, separator, uptimeStr, separator, levelStr)
	return styles.header.Render(left)
}

func (m model) renderSummaryBar(styles styles) string {
	total, running, queued, pass, fail, warn := m.scheduler.GetCounts()

	totalStr := styles.summaryText.Render(fmt.Sprintf("Total: %d", total))
	runningStr := styles.summaryText.Render(fmt.Sprintf("Running: %d", running))
	queuedStr := styles.summaryText.Render(fmt.Sprintf("Queued: %d", queued))
	passStr := styles.summaryText.Render(fmt.Sprintf("Pass: %d", pass))
	failStr := styles.summaryText.Render(fmt.Sprintf("Fail: %d", fail))
	warnStr := styles.summaryText.Render(fmt.Sprintf("Warn: %d", warn))

	separator := styles.summaryText.Render(" | ")
	content := lipgloss.JoinHorizontal(lipgloss.Top, totalStr, separator, runningStr, separator, queuedStr, separator, passStr, separator, failStr, separator, warnStr)

	return styles.summaryBar.Render(content)
}

func (m model) renderCheckList(styles styles) string {
	checks := m.scheduler.GetAllChecks()

	if len(checks) == 0 {
		return styles.empty.Render("No checks configured")
	}

	var names []string
	for name := range checks {
		names = append(names, name)
	}
	sort.Strings(names)

	availableWidth := styles.width - 2
	nameWidth := m.calculateNameWidth(names, availableWidth)
	var rows []string
	now := time.Now()

	for _, name := range names {
		check := checks[name]
		rows = append(rows, m.formatCheckRow(name, nameWidth, check, now, styles))
	}

	totalRows := len(rows)
	availableRows := m.getAvailableCheckRows()
	maxRows := availableRows - 1
	if totalRows > maxRows {
		maxRows--
	}
	if maxRows > m.maxCheckRows {
		maxRows = m.maxCheckRows
	}
	if maxRows < 1 {
		maxRows = 1
	}

	displayStart := 0
	if totalRows > maxRows {
		scrollWindow := totalRows - maxRows + 1
		seconds := int(time.Since(m.uptime).Seconds())
		displayStart = seconds % scrollWindow
	}
	displayEnd := min(displayStart+maxRows, totalRows)

	var displayRows []string
	displayRows = append(displayRows, m.renderCheckHeader(nameWidth, styles))
	displayRows = append(displayRows, m.renderCheckDivider(styles))
	for i := displayStart; i < displayEnd; i++ {
		displayRows = append(displayRows, rows[i])
	}

	content := strings.Join(displayRows, "\n")

	if totalRows > maxRows {
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d", displayStart+1, displayEnd, totalRows)
		content += "\n" + styles.scrollInfo.Render(scrollInfo)
	}

	return styles.checkList.Render(content)
}

func (m model) renderCheckHeader(nameWidth int, styles styles) string {
	return styles.columnHeader.Render(m.formatCheckRow("Check", nameWidth, nil, time.Now(), styles))
}

func (m model) renderCheckDivider(styles styles) string {
	width := styles.width - 2
	if width < 1 {
		width = 1
	}
	return styles.divider.Render(strings.Repeat("─", width))
}

func (m model) formatCheckRow(name string, nameWidth int, check *scheduler.ScheduledCheck, now time.Time, styles styles) string {
	var status string
	if check == nil {
		status = "Status"
	} else if check.Running {
		status = styles.colorRunning.Render("⟳")
	} else if check.IsQueued {
		status = styles.colorQueued.Render("⏳")
	} else {
		status = styles.colorIdle.Render("•")
	}

	var result string
	if check == nil {
		result = "Last"
	} else {
		switch check.LastStatus {
		case "pass":
			result = styles.colorPass.Render("✓")
		case "fail":
			result = styles.colorFail.Render("✗")
		case "warn":
			result = styles.colorWarn.Render("⚠")
		default:
			result = styles.colorUnknown.Render("?")
		}
	}

	var lastRun string
	if check == nil {
		lastRun = "Last Run"
	} else if check.LastRun != nil {
		lastRun = formatRelativeTime(now.Sub(*check.LastRun))
	} else {
		lastRun = "never"
	}

	var nextRun string
	if check == nil {
		nextRun = "Next Run"
	} else if check.NextRun.After(now) {
		nextRun = fmt.Sprintf("in %s", formatRelativeTime(check.NextRun.Sub(now)))
	} else {
		nextRun = "due"
	}

	statusWidth := 6
	resultWidth := 5
	lastWidth := 12
	nextWidth := 12
	availableWidth := styles.width - 2

	nameCell := padRight(truncate(name, nameWidth), nameWidth)
	statusCell := padRight(status, statusWidth)
	resultCell := padRight(result, resultWidth)
	lastCell := padRight(lastRun, lastWidth)
	nextCell := padRight(nextRun, nextWidth)

	row := fmt.Sprintf("%s  %s  %s  %s  %s",
		nameCell, statusCell, resultCell, lastCell, nextCell)
	return padRight(row, availableWidth)
}

func (m model) renderFooter(styles styles) string {
	help := styles.footerText.Render("Ctrl+C to exit")
	refresh := styles.footerText.Render("Refresh: 1s")
	separator := styles.footerText.Render("    ")

	return styles.footer.Render(lipgloss.JoinHorizontal(lipgloss.Top, help, separator, refresh))
}

func formatRelativeTime(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

func truncate(s string, maxLen int) string {
	if runewidth.StringWidth(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return runewidth.Truncate(s, maxLen, "")
	}
	return runewidth.Truncate(s, maxLen, "...")
}

func padRight(s string, width int) string {
	if width <= 0 {
		return ""
	}
	current := lipgloss.Width(s)
	if current >= width {
		return s
	}
	return s + strings.Repeat(" ", width-current)
}

func (m model) calculateNameWidth(names []string, availableWidth int) int {
	statusWidth := 6
	resultWidth := 5
	lastWidth := 12
	nextWidth := 12
	minNameWidth := 10
	maxNameWidth := 32

	reserved := statusWidth + resultWidth + lastWidth + nextWidth + 8
	nameWidth := availableWidth - reserved
	if nameWidth < minNameWidth {
		nameWidth = minNameWidth
	}
	if nameWidth > maxNameWidth {
		nameWidth = maxNameWidth
	}

	maxSeen := runewidth.StringWidth("Check")
	for _, name := range names {
		width := runewidth.StringWidth(name)
		if width > maxSeen {
			maxSeen = width
		}
	}
	if maxSeen < minNameWidth {
		maxSeen = minNameWidth
	}
	if maxSeen > nameWidth {
		return nameWidth
	}
	return maxSeen
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
