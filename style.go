package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// red
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// grey
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// green
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#84D2C5"))
	// aqa
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	// dark grey
	barStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#42855B"))

	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

// Add a purple, rectangular border
var style = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

// Set a rounded, yellow-on-purple border to the top and left
var anotherStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("228")).
	BorderBackground(lipgloss.Color("63")).
	BorderTop(true).
	BorderLeft(true)

// Make your own border
var myCuteBorder = lipgloss.Border{
	Top:         "._.:*:",
	Bottom:      "._.:*:",
	Left:        "|*",
	Right:       "|*",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}
