package editor

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps          = 60
	springFreq   = 6.0  // oscillations per second — higher = snappier
	springDamp   = 0.8  // 0=bouncy, 1=no bounce; 0.8 feels crisp
	panelSpringF = 4.0  // slower spring for preview panel slide-in
	panelSpringD = 0.9
)

// tickMsg drives the animation loop
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Field represents one editable theme property
type Field struct {
	Label string
	Key   string
	Value string
}

// Model is the Bubbletea model for the theme editor
type Model struct {
	theme     *config.Theme
	themePath string
	fields    []Field
	cursor    int
	editing   bool
	input     textinput.Model
	saved     bool
	message   string
	width     int
	height    int

	// ── harmonica spring state ──────────────────────────
	// cursor spring — animates the selector row highlight
	cursorSpring harmonica.Spring
	cursorPos    float64 // current animated position (fractional row index)
	cursorVel    float64 // current velocity

	// panel spring — animates preview panel sliding in on open
	panelSpring  harmonica.Spring
	panelOffset  float64 // 0 = fully visible, 1 = off-screen right
	panelVel     float64
	panelReady   bool // true once panel has animated in
}

// NewModel creates the editor model from a theme
func NewModel(t *config.Theme, path string) Model {
	ti := textinput.New()
	ti.CharLimit = 64

	fields := buildFields(t)

	return Model{
		theme:        t,
		themePath:    path,
		fields:       fields,
		input:        ti,
		cursorSpring: harmonica.NewSpring(harmonica.FPS(fps), springFreq, springDamp),
		cursorPos:    0,
		cursorVel:    0,
		panelSpring:  harmonica.NewSpring(harmonica.FPS(fps), panelSpringF, panelSpringD),
		panelOffset:  1, // start off-screen
		panelVel:     0,
		panelReady:   false,
	}
}

func buildFields(t *config.Theme) []Field {
	return []Field{
		{Label: "Name", Key: "meta.name", Value: t.Meta.Name},
		{Label: "Author", Key: "meta.author", Value: t.Meta.Author},
		{Label: "Description", Key: "meta.description", Value: t.Meta.Description},
		{Label: "Primary Color", Key: "colors.primary", Value: t.Colors.Primary},
		{Label: "Secondary Color", Key: "colors.secondary", Value: t.Colors.Secondary},
		{Label: "Background", Key: "colors.background", Value: t.Colors.Background},
		{Label: "Foreground", Key: "colors.foreground", Value: t.Colors.Foreground},
		{Label: "Accent Color", Key: "colors.accent", Value: t.Colors.Accent},
		{Label: "Error Color", Key: "colors.error", Value: t.Colors.Error},
		{Label: "Success Color", Key: "colors.success", Value: t.Colors.Success},
		{Label: "Warning Color", Key: "colors.warning", Value: t.Colors.Warning},
		{Label: "Muted Color", Key: "colors.muted", Value: t.Colors.Muted},
		{Label: "Prompt Symbol", Key: "prompt.symbol", Value: t.Prompt.Symbol},
		{Label: "Prompt Style", Key: "prompt.style", Value: t.Prompt.Style},
		{Label: "Prompt Format", Key: "prompt.format", Value: t.Prompt.Format},
		{Label: "Loader Speed (ms)", Key: "loader.interval_ms", Value: fmt.Sprintf("%d", t.Loader.IntervalMs)},
		{Label: "Banner Text", Key: "banner.text", Value: t.Banner.Text},
		{Label: "Banner Effect", Key: "graphics.effects.banner", Value: t.Graphics.Effects.Banner},
		{Label: "Gradient From", Key: "graphics.gradient.from", Value: t.Graphics.Gradient.From},
		{Label: "Gradient To", Key: "graphics.gradient.to", Value: t.Graphics.Gradient.To},
		{Label: "Divider Style", Key: "graphics.divider.style", Value: t.Graphics.Divider.Style},
		{Label: "Cursor Style", Key: "cursor.style", Value: t.Cursor.Style},
		{Label: "Progress Filled", Key: "progress_bar.filled", Value: t.ProgressBar.Filled},
		{Label: "Progress Empty", Key: "progress_bar.empty", Value: t.ProgressBar.Empty},
	}
}

// Init starts the harmonica animation loop (cursor spring + panel
// slide-in) as soon as the editor opens. Satisfies tea.Model.
func (m Model) Init() tea.Cmd {
	// start animation loop immediately
	return tick()
}

// Update handles key input and animation ticks, advancing the cursor and
// panel springs each frame. Satisfies tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		// advance cursor spring toward target (float of cursor index)
		target := float64(m.cursor)
		m.cursorPos, m.cursorVel = m.cursorSpring.Update(m.cursorPos, m.cursorVel, target)

		// advance panel spring toward 0 (fully visible)
		m.panelOffset, m.panelVel = m.panelSpring.Update(m.panelOffset, m.panelVel, 0)
		if m.panelOffset < 0.01 {
			m.panelOffset = 0
			m.panelReady = true
		}

		return m, tick()

	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				m.fields[m.cursor].Value = m.input.Value()
				applyField(m.theme, m.fields[m.cursor])
				m.editing = false
				m.input.Blur()
				m.message = fmt.Sprintf("Updated: %s", m.fields[m.cursor].Label)
			case "esc":
				m.editing = false
				m.input.Blur()
				m.message = "Edit cancelled"
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					// spring will chase new target on next tick
				}
			case "down", "j":
				if m.cursor < len(m.fields)-1 {
					m.cursor++
				}
			case "enter":
				m.editing = true
				m.input.SetValue(m.fields[m.cursor].Value)
				m.input.Focus()
				m.message = fmt.Sprintf("Editing: %s", m.fields[m.cursor].Label)
			case "s":
				if err := saveTheme(m.theme, m.themePath); err != nil {
					m.message = "✗ Save failed: " + err.Error()
				} else {
					m.saved = true
					m.message = "✓ Theme saved to " + filepath.Base(m.themePath)
				}
			}
		}
	}
	return m, nil
}

// View renders the two-panel editor layout: the scrollable field list on
// the left (with the spring-animated cursor highlight) and the live theme
// preview on the right (with the spring-animated slide-in). Satisfies
// tea.Model.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading editor..."
	}

	primary := m.theme.Colors.Primary
	accent := m.theme.Colors.Accent
	muted := m.theme.Colors.Muted

	// ── Styles ─────────────────────────────────────────
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primary)).
		Bold(true).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(accent)).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(muted))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primary)).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(accent))

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(muted))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(muted)).
		Italic(true)

	// ── Title ──────────────────────────────────────────
	title := titleStyle.Render(fmt.Sprintf("cmdX Theme Editor — %s", m.theme.Meta.Name))

	// ── Fields panel ───────────────────────────────────
	// animated cursor row: round the spring position for highlight
	animatedRow := int(math.Round(m.cursorPos))

	var fieldLines []string
	visibleStart := 0
	visibleCount := 18

	if m.cursor >= visibleStart+visibleCount {
		visibleStart = m.cursor - visibleCount + 1
	}

	for i := visibleStart; i < len(m.fields) && i < visibleStart+visibleCount; i++ {
		f := m.fields[i]

		// use animated row for the highlight so it springs smoothly
		isHighlighted := i == animatedRow
		isCursor := i == m.cursor

		prefix := "  "
		if isHighlighted {
			prefix = "▶ "
		}

		var line string
		if isCursor && m.editing {
			line = fmt.Sprintf("%s%s  %s",
				prefix,
				labelStyle.Render(f.Label+":"),
				m.input.View(),
			)
		} else if isHighlighted {
			line = selectedStyle.Render(fmt.Sprintf("%s%-22s  %s", prefix, f.Label+":", f.Value))
		} else {
			line = fmt.Sprintf("%s%s  %s",
				prefix,
				normalStyle.Render(fmt.Sprintf("%-22s", f.Label+":")),
				mutedStyle.Render(f.Value),
			)
		}
		fieldLines = append(fieldLines, line)
	}

	fieldsPanel := strings.Join(fieldLines, "\n")

	// ── Preview panel ───────────────────────────────────
	// apply panel spring offset as right-padding to simulate slide-in
	preview := m.buildPreviewPanel(labelStyle, valueStyle, mutedStyle)

	rightWidth := 35
	// clamp panelOffset so padding doesn't go negative
	panelPad := int(math.Round(m.panelOffset * float64(rightWidth)))
	if panelPad < 0 {
		panelPad = 0
	}
	if panelPad > rightWidth {
		panelPad = rightWidth
	}

	// ── Layout ─────────────────────────────────────────
	leftWidth := 45

	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(primary)).
		Width(leftWidth).
		Height(22).
		Padding(0, 1).
		Render(fieldsPanel)

	rightBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(accent)).
		Width(rightWidth).
		Height(22).
		Padding(0, 1).
		MarginLeft(panelPad). // spring-driven slide-in margin
		Render(preview)

	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, "  ", rightBox)

	// ── Message bar ────────────────────────────────────
	statusMsg := helpStyle.Render(m.message)
	if m.message == "" {
		statusMsg = helpStyle.Render("[↑↓] Navigate  [Enter] Edit  [S] Save  [Q] Quit")
	}

	return fmt.Sprintf("%s\n%s\n%s", title, columns, statusMsg)
}

func (m Model) buildPreviewPanel(labelStyle, valueStyle, mutedStyle lipgloss.Style) string {
	t := m.theme
	var lines []string

	// colors
	lines = append(lines, labelStyle.Render("Colors:"))
	colorPairs := [][2]string{
		{"primary", t.Colors.Primary},
		{"accent", t.Colors.Accent},
		{"success", t.Colors.Success},
		{"error", t.Colors.Error},
	}
	for _, pair := range colorPairs {
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(pair[1])).
			Render("  ")
		lines = append(lines, fmt.Sprintf("  %s %s", swatch,
			valueStyle.Render(pair[0])))
	}

	lines = append(lines, "")

	// prompt preview
	lines = append(lines, labelStyle.Render("Prompt:"))
	dirStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Colors.Accent)).Bold(true)
	symbolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Colors.Primary)).Bold(true)
	lines = append(lines, fmt.Sprintf("  %s %s",
		dirStyle.Render("~/projects"),
		symbolStyle.Render(t.Prompt.Symbol),
	))

	lines = append(lines, "")

	// loader frames
	lines = append(lines, labelStyle.Render("Loader:"))
	frames := ""
	for _, f := range t.Loader.Frames {
		frames += f + " "
	}
	lines = append(lines, "  "+mutedStyle.Render(frames))

	lines = append(lines, "")

	// progress bar
	lines = append(lines, labelStyle.Render("Progress:"))
	filled := strings.Repeat(t.ProgressBar.Filled, 12)
	empty := strings.Repeat(t.ProgressBar.Empty, 8)
	barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Colors.Primary))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Colors.Muted))
	lines = append(lines, fmt.Sprintf("  [%s%s]",
		barStyle.Render(filled),
		emptyStyle.Render(empty),
	))

	lines = append(lines, "")

	// banner
	lines = append(lines, labelStyle.Render("Banner:"))
	bannerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Colors.Primary)).Bold(true)
	lines = append(lines, "  "+bannerStyle.Render(t.Banner.Text))

	return strings.Join(lines, "\n")
}

// applyField writes a field value back into the theme struct
func applyField(t *config.Theme, f Field) {
	switch f.Key {
	case "meta.name":
		t.Meta.Name = f.Value
	case "meta.author":
		t.Meta.Author = f.Value
	case "meta.description":
		t.Meta.Description = f.Value
	case "colors.primary":
		t.Colors.Primary = f.Value
	case "colors.secondary":
		t.Colors.Secondary = f.Value
	case "colors.background":
		t.Colors.Background = f.Value
	case "colors.foreground":
		t.Colors.Foreground = f.Value
	case "colors.accent":
		t.Colors.Accent = f.Value
	case "colors.error":
		t.Colors.Error = f.Value
	case "colors.success":
		t.Colors.Success = f.Value
	case "colors.warning":
		t.Colors.Warning = f.Value
	case "colors.muted":
		t.Colors.Muted = f.Value
	case "prompt.symbol":
		t.Prompt.Symbol = f.Value
	case "prompt.style":
		t.Prompt.Style = f.Value
	case "prompt.format":
		t.Prompt.Format = f.Value
	case "banner.text":
		t.Banner.Text = f.Value
	case "cursor.style":
		t.Cursor.Style = f.Value
	case "progress_bar.filled":
		t.ProgressBar.Filled = f.Value
	case "progress_bar.empty":
		t.ProgressBar.Empty = f.Value
	case "graphics.effects.banner":
		t.Graphics.Effects.Banner = f.Value
	case "graphics.gradient.from":
		t.Graphics.Gradient.From = f.Value
	case "graphics.gradient.to":
		t.Graphics.Gradient.To = f.Value
	case "graphics.divider.style":
		t.Graphics.Divider.Style = f.Value
	}
}

// saveTheme writes the updated theme back to JSON
func saveTheme(t *config.Theme, path string) error {
	data, err := marshalTheme(t)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
