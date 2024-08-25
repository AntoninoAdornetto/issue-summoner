package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Selection struct {
	Options map[int]bool
}

type Step struct {
	StepName string
	Options  []Item
	Headers  string
	Field    string
}

type Item struct {
	ID          int
	Title, Desc string
}

func (s *Selection) OnSelect(option int, value bool) {
	s.Options[option] = value
}

type model struct {
	cursor   int
	options  []Item
	selected map[int]struct{}
	choices  *Selection
	header   string
	exit     *bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitMultiSelect(
	options []Item,
	selection *Selection,
	header string,
	program *bool,
) model {
	return model{
		options:  options,
		selected: make(map[int]struct{}),
		choices:  selection,
		header:   AccentTextStyle.Render(header),
		exit:     program,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			*m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "y":
			for selectedKey := range m.selected {
				m.choices.OnSelect(m.options[selectedKey].ID, true)
				m.cursor = selectedKey
			}
			*m.exit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(m.header + "\n\n")

	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = SuccessTextStyle.Render(">")
			option.Title = PrimaryTextStyle.Render(option.Title)
			option.Desc = PrimaryTextStyle.Render(option.Desc)
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = SecondaryTextStyle.Render("*")
		}

		title := DimTextStyle.Render(option.Title)
		description := DimTextStyle.Render(option.Desc)

		s.WriteString(fmt.Sprintf("%s [%s] %s\n%s\n\n", cursor, checked, title, description))
	}

	s.WriteString(AccentTextStyle.Render("\u2191 "))
	s.WriteString("or ")
	s.WriteString(AccentTextStyle.Render("k "))
	s.WriteString("= move up list")

	s.WriteString(AccentTextStyle.Render("\n\u2193 "))
	s.WriteString("or ")
	s.WriteString(AccentTextStyle.Render("j "))
	s.WriteString("= move down list")

	s.WriteString(AccentTextStyle.Render("\nspace "))
	s.WriteString("= select/deselect")

	s.WriteString(AccentTextStyle.Render("\ny "))
	s.WriteString("= confirm choices\n")

	return s.String()
}
