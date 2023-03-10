package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println("ToggleCheck panicked, may caused by slice index out of bounce (length==1)")
					}
				}()
				dao.ToggleCheck(m.potatoes.Choices[m.potatoes.Cursor].ID)
			}()

		// append normal
		case key.Matches(msg, m.keymap.Append):
			return InitModelAddEntry(dao.NORMAL), nil

		// append daily
		case key.Matches(msg, m.keymap.AppendDaily):
			return InitModelAddEntry(dao.DAILY), nil

		// delete
		case key.Matches(msg, m.keymap.Delete):
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println("DeleteEntry panicked, may caused by slice index out of bounce (length==1)")
					}
				}()
				dao.DeleteEntry(m.potatoes.Choices[m.potatoes.Cursor].ID)
			}()
			return initModelDashboard(), nil

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
	// s := focusedStyle.Render("Welcome to ???? Potatoes ???? -- A personal flavored TODO Applet\nWhat will we do?\n\n")
	// s := style.Render("Welcome to ???? Potatoes ???? -- A personal flavored TODO Applet\nWhat will we do?\n\n")

	// Vertically join two paragraphs along their center axes
	s := focusedStyle.Render(lipgloss.JoinVertical(lipgloss.Left, "Welcome to ???? Potatoes ????", strings.Repeat(" ", 15)+"-- A personal flavored TODO Applet", "What will we do?"))
	// Center a paragraph horizontally in a space 80 cells wide. The height of
	// the block returned will be as tall as the input paragraph.
	s = lipgloss.PlaceHorizontal(50, lipgloss.Left, s)

	d := noStyle.Render(lipgloss.PlaceHorizontal(50, lipgloss.Center, "daily", lipgloss.WithWhitespaceChars("="))) + "\n"
	n := noStyle.Render(lipgloss.PlaceHorizontal(50, lipgloss.Center, "normal", lipgloss.WithWhitespaceChars("="))) + "\n"

	// daily counter / normal counter
	dc := 0
	nc := 0

	for i, choice := range m.potatoes.Choices {
		cursor := " "
		if i == m.potatoes.Cursor {
			cursor = focusedStyle.Render(">")
		}

		checked := " "
		if _, ok := m.potatoes.Checked[i]; ok {
			checked = "x"

			if i == m.potatoes.Cursor {
				checked = focusedStyle.Render("x")
			}

		}

		if choice.Type == dao.DAILY {
			d += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Entry)
			dc++
		} else if choice.Type == dao.NORMAL {
			n += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Entry)
			nc++
		}
	}

	// extra hint
	if dc < 1 {
		d += titleStyle.Render("\nNo pending daily chores yet, press `A` to add one! ????\n")
	}
	if nc < 1 {
		n += titleStyle.Render("\nAll good, everything is going so well! ????\n")
	}

	helpView := m.help.View(m.keymap)
	height := 5 - strings.Count(s, "\n") - strings.Count(helpView, "\n")

	s = lipgloss.JoinVertical(lipgloss.Left, s, d, n, strings.Repeat("\n", height)+helpView)

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
	// dao.InitMySQL()
	dao.InitSQLiet()

	if err := tea.NewProgram(initModelDashboard()).Start(); err != nil {
		fmt.Println("Oooouch, something bad...")
	}
}
