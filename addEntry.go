package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	addentry "rey.com/charm/potatoes/addEntry"
	"rey.com/charm/potatoes/dao"
)

type (
	errMsg error
)

type ModelAddEntry struct {
	textInput textinput.Model
	err       error

	keymap addentry.KeyMap
	help   help.Model

	potatoType dao.PotatoType
}

func InitModelAddEntry(t dao.PotatoType) tea.Model {
	ti := textinput.New()
	if t == dao.NORMAL {
		ti.Placeholder = "The next thing to be done... 🤔"
	} else if t == dao.DAILY {
		ti.Placeholder = "Think about \"21-Day Rule\" 🌞"
	}
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.PromptStyle = focusedStyle

	return ModelAddEntry{
		textInput: ti,
		err:       nil,

		keymap: addentry.DefaultKeyMap,
		help:   help.New(),

		potatoType: t,
	}
}

func (m ModelAddEntry) Init() tea.Cmd {
	return textinput.Blink
}

func (m ModelAddEntry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:

		switch {
		// quit and go back to potatoes
		case key.Matches(msg, m.keymap.Quit):
			return initModelDashboard(), nil

		// confirm
		case key.Matches(msg, m.keymap.Confirm):
			input := strings.TrimSpace(m.textInput.Value())

			if len(input) < 1 {
				return initModelDashboard(), nil
			}

			if m.potatoType == dao.DAILY {
				return AppendModelDashboard(input, dao.DAILY), nil
			} else if m.potatoType == dao.NORMAL {
				return AppendModelDashboard(input, dao.NORMAL), nil
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m ModelAddEntry) View() string {
	var prompt string
	if m.potatoType == dao.NORMAL {
		prompt = "The next thing to be done... 🤔"
	} else if m.potatoType == dao.DAILY {
		prompt = "Create a daily chore 🌞"
	}

	s := fmt.Sprintf(
		"%s\n\n%s",
		prompt,
		m.textInput.View(),
	) + "\n"

	helpView := m.help.View(m.keymap)
	height := 5 - strings.Count(s, "\n") - strings.Count(helpView, "\n")

	s += strings.Repeat("\n", height) + helpView

	return s
}
