package dashboard

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type KeyMap struct {
	Up   key.Binding
	Down key.Binding
	// Left  key.Binding
	// Right key.Binding
	Select key.Binding
	Append key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},       // first column
		{k.Select, k.Append}, // second column
		{k.Help, k.Quit},
	}
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	//Left: key.NewBinding(
	//key.WithKeys("left", "h"),
	//key.WithHelp("←/h", "move left"),
	//),
	//Right: key.NewBinding(
	//key.WithKeys("right", "l"),
	//key.WithHelp("→/l", "move right"),
	//),
	Select: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("enter", "select a item"),
	),
	Append: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "append a new entry"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var ColemakKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "u"),
		key.WithHelp("↑/u", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "e"),
		key.WithHelp("↓/e", "move down"),
	),
	//Left: key.NewBinding(
	//key.WithKeys("left", "n"),
	//key.WithHelp("←/n", "move left"),
	//),
	//Right: key.NewBinding(
	//key.WithKeys("right", "i"),
	//key.WithHelp("→/i", "move right"),
	//),
	Select: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("enter", "select a item"),
	),
	Append: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "append a new entry"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
