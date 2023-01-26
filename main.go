package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"rey.com/charm/potatoes/bubblesx"
)

type model struct {
	potatoes Potatoes
	keymap   bubblesx.KeyMap
	help     help.Model
}

type Potatoes struct {
	Choices  []string
	Cursor   int
	Selected map[int]struct{}
}

func initModel() *model {
	return &model{
		potatoes: Potatoes{
			Choices:  []string{"Buy carrots", "Buy chocolates", "Buy milk"},
			Selected: make(map[int]struct{}),
		},
		keymap: bubblesx.ColemakKeyMap,
		help:   help.New(),
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *model) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {

		// quit
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		// up
		case key.Matches(msg, m.keymap.Up):
			if m.potatoes.Cursor > 0 {
				m.potatoes.Cursor--
			}

		// down
		case key.Matches(msg, m.keymap.Down):
			if m.potatoes.Cursor < len(m.potatoes.Choices)-1 {
				m.potatoes.Cursor++
			}

		// select
		case key.Matches(msg, m.keymap.Select):
			_, ok := m.potatoes.Selected[m.potatoes.Cursor]
			if ok {
				delete(m.potatoes.Selected, m.potatoes.Cursor)
			} else {
				m.potatoes.Selected[m.potatoes.Cursor] = struct{}{}
			}

		// help
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *model) View() string {
	s := "What will we buy?\n\n"

	for i, choice := range m.potatoes.Choices {
		cursor := " "
		if i == m.potatoes.Cursor {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.potatoes.Selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	helpView := m.help.View(m.keymap)
	height := 8 - strings.Count(s, "\n") - strings.Count(helpView, "\n")

	s += strings.Repeat("\n", height) + helpView

	return s
}

func main() {
	os.Setenv("HELP_DEBUG", "enable")
	if os.Getenv("HELP_DEBUG") != "" {
		if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}
	if err := tea.NewProgram(initModel()).Start(); err != nil {
		fmt.Println("Oooouch, something bad...")
	}
}
