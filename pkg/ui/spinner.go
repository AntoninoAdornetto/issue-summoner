package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type spinnerModel struct {
	spinner spinner.Model
	exit    bool
	err     error
	msg     string
}

func InitSpinner(str string) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(AccentTextStyle.GetForeground())
	return spinnerModel{spinner: s, msg: str}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.exit = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {

	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("%s %s\n\n", m.spinner.View(), m.msg)
	if m.exit {
		return str + "\n"
	}
	return str
}
