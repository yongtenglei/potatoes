package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"rey.com/charm/potatoes/dao"
	"rey.com/charm/potatoes/dashboard"
)

type ModelDashboard struct {
	potatoes Potatoes
	keymap   dashboard.KeyMap
	help     help.Model
}

type Potatoes struct {
	Choices []dao.Potato
	Cursor  int
	Checked map[int]struct{}
}

func initModelDashboard() tea.Model {
	potatoes, err := dao.LoadPotatoes()
	if err != nil {
		fmt.Println("Oooouch, something bad...")
		os.Exit(1)
	}

	checked := make(map[int]struct{})
	for i := 0; i < len(potatoes); i++ {
		if potatoes[i].Checked {
			checked[i] = struct{}{}
		}
	}

	return &ModelDashboard{
		potatoes: Potatoes{
			Choices: potatoes,
			Checked: checked,
		},
		keymap: dashboard.ColemakKeyMap,
		help:   help.New(),
	}
}

func AppendModelDashboard(entry string, potatoType dao.PotatoType) tea.Model {
	if err := dao.AddEntry(entry, potatoType); err != nil {
		log.Println(err)
	}

	potatoes, err := dao.LoadPotatoes()
	if err != nil {
		log.Println("Oooouch, something bad...")
		os.Exit(1)
	}

	checked := make(map[int]struct{})
	for i := 0; i < len(potatoes); i++ {
		if potatoes[i].Checked {
			checked[i] = struct{}{}
		}
	}

	return &ModelDashboard{
		potatoes: Potatoes{
			Choices: potatoes,
			Checked: checked,
		},
		keymap: dashboard.ColemakKeyMap,
		help:   help.New(),
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *ModelDashboard) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *ModelDashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		// cycle cursor (up and down)
		// up
		case key.Matches(msg, m.keymap.Up):
			if m.potatoes.Cursor > 0 {
				m.potatoes.Cursor--
			} else {
				m.potatoes.Cursor = len(m.potatoes.Choices) - 1
			}

		// down
		case key.Matches(msg, m.keymap.Down):
			if m.potatoes.Cursor < len(m.potatoes.Choices)-1 {
				m.potatoes.Cursor++
			} else {
				m.potatoes.Cursor = 0
			}

		// select
		case key.Matches(msg, m.keymap.Select):
			_, ok := m.potatoes.Checked[m.potatoes.Cursor]
			if ok {
				delete(m.potatoes.Checked, m.potatoes.Cursor)
			} else {
				m.potatoes.Checked[m.potatoes.Cursor] = struct{}{}
			}
			// PERF: Manipulate database at any interactive point. Could be some kind lazy interaction?
			// WARN: May caused inconsistency.
			go dao.ToggleCheck(m.potatoes.Choices[m.potatoes.Cursor].ID)

		// append normal
		case key.Matches(msg, m.keymap.Append):
			return InitModelAddEntry(dao.NORMAL), nil

		// append daily
		case key.Matches(msg, m.keymap.AppendDaily):
			return InitModelAddEntry(dao.DAILY), nil

		// help
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *ModelDashboard) View() string {
	s := "What will we buy?\n\n"
	d := strings.Repeat("=", 22) + "daily" + strings.Repeat("=", 22) + "\n"
	n := strings.Repeat("=", 22) + "normal" + strings.Repeat("=", 21) + "\n"

	for i, choice := range m.potatoes.Choices {
		cursor := " "
		if i == m.potatoes.Cursor {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.potatoes.Checked[i]; ok {
			checked = "x"
		}

		if choice.Type == dao.DAILY {
			d += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Entry)
		} else if choice.Type == dao.NORMAL {
			n += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Entry)
		}
	}

	helpView := m.help.View(m.keymap)
	height := 8 - strings.Count(s, "\n") - strings.Count(helpView, "\n")

	s = s + d + n
	s += strings.Repeat("\n", height) + helpView

	return s
}

func main() {
	// always enable debug mode for now
	os.Setenv("HELP_DEBUG", "enable")

	if os.Getenv("HELP_DEBUG") != "" {
		if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}

	// connect to SQL
	dao.InitMySQL()

	if err := tea.NewProgram(initModelDashboard()).Start(); err != nil {
		fmt.Println("Oooouch, something bad...")
	}
}
