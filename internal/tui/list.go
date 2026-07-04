// Package tui provides reusable Bubble Tea TUI components for cmdX.
package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ThemeItem is a list.Item representing one theme entry.
type ThemeItem struct {
	Name   string
	Desc   string
	Author string
}

// Title, Description, and FilterValue satisfy bubbles/list.Item.
func (t ThemeItem) Title() string       { return t.Name }
func (t ThemeItem) Description() string { return t.Desc }
func (t ThemeItem) FilterValue() string { return t.Name }

// themeDelegate controls how each item is rendered.
type themeDelegate struct {
	primary string
	accent  string
	muted   string
}

// Height, Spacing, Update, and Render satisfy bubbles/list.ItemDelegate.
func (d themeDelegate) Height() int                             { return 2 }
func (d themeDelegate) Spacing() int                            { return 1 }
func (d themeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d themeDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(ThemeItem)
	if !ok {
		return
	}

	selected := index == m.Index()

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(d.muted))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(d.muted)).Faint(true)
	prefix := "  "

	if selected {
		nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(d.accent)).Bold(true)
		descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(d.primary))
		prefix = "▶ "
	}

	name := fmt.Sprintf("%s%s", prefix, nameStyle.Render(item.Name))
	desc := fmt.Sprintf("  %s", descStyle.Render(item.Desc))
	fmt.Fprintln(w, name)
	fmt.Fprintln(w, desc)
}

// ThemeListModel is a full Bubble Tea model for interactive theme selection.
type ThemeListModel struct {
	list     list.Model
	chosen   string
	quitting bool
	primary  string
	accent   string
}

// NewThemeList builds an interactive theme picker from a slice of ThemeItems.
func NewThemeList(items []ThemeItem, primary, accent, muted string) ThemeListModel {
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = it
	}

	delegate := themeDelegate{primary: primary, accent: accent, muted: muted}
	l := list.New(listItems, delegate, 50, 20)
	l.Title = "Available Themes"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(primary)).
		Bold(true).
		Padding(0, 1)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color(accent))
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color(accent))
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	return ThemeListModel{list: l, primary: primary, accent: accent}
}

// Init, Update, and View satisfy tea.Model.
func (m ThemeListModel) Init() tea.Cmd { return nil }

func (m ThemeListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(ThemeItem); ok {
				m.chosen = item.Name
				return m, tea.Quit
			}
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ThemeListModel) View() string {
	if m.quitting {
		return ""
	}
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.primary)).
		Padding(1, 2)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true).
		Render("[↑↓] Navigate  [/] Filter  [Enter] Select  [Q] Quit")

	return border.Render(m.list.View()) + "\n" + help
}

// Chosen returns the selected theme name, or empty string if cancelled.
func (m ThemeListModel) Chosen() string { return m.chosen }

// RunThemeList launches the interactive theme picker and returns the chosen name.
func RunThemeList(items []ThemeItem, primary, accent, muted string) (string, error) {
	m := NewThemeList(items, primary, accent, muted)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	return result.(ThemeListModel).Chosen(), nil
}

// RegistryItem is a list.Item for registry theme entries.
type RegistryItem struct {
	Name    string
	Author  string
	Desc    string
	Version string
	Tags    []string
}

// Title, Description, and FilterValue satisfy bubbles/list.Item.
func (r RegistryItem) Title() string { return r.Name }
func (r RegistryItem) Description() string {
	tags := strings.Join(r.Tags, ", ")
	if tags == "" {
		return fmt.Sprintf("by %s — %s", r.Author, r.Desc)
	}
	return fmt.Sprintf("by %s — %s  [%s]", r.Author, r.Desc, tags)
}
func (r RegistryItem) FilterValue() string { return r.Name + " " + strings.Join(r.Tags, " ") }
