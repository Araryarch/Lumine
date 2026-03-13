package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TextInput struct {
	Prompt      string
	Placeholder string
	Value       string
	CursorPos   int
	Width       int
	Focused     bool
}

func NewTextInput(prompt, placeholder string, width int) TextInput {
	return TextInput{
		Prompt:      prompt,
		Placeholder: placeholder,
		Width:       width,
		CursorPos:   0,
		Focused:     false,
	}
}

func (t *TextInput) InsertChar(ch rune) {
	if t.CursorPos > len(t.Value) {
		t.CursorPos = len(t.Value)
	}

	before := t.Value[:t.CursorPos]
	after := t.Value[t.CursorPos:]
	t.Value = before + string(ch) + after
	t.CursorPos++
}

func (t *TextInput) DeleteChar() {
	if t.CursorPos > 0 {
		before := t.Value[:t.CursorPos-1]
		after := t.Value[t.CursorPos:]
		t.Value = before + after
		t.CursorPos--
	}
}

func (t *TextInput) MoveCursorLeft() {
	if t.CursorPos > 0 {
		t.CursorPos--
	}
}

func (t *TextInput) MoveCursorRight() {
	if t.CursorPos < len(t.Value) {
		t.CursorPos++
	}
}

func (t *TextInput) MoveCursorStart() {
	t.CursorPos = 0
}

func (t *TextInput) MoveCursorEnd() {
	t.CursorPos = len(t.Value)
}

func (t *TextInput) Clear() {
	t.Value = ""
	t.CursorPos = 0
}

func (t TextInput) Render() string {
	var b strings.Builder

	promptStyle := lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true)
	b.WriteString(promptStyle.Render(t.Prompt))
	b.WriteString(" ")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(t.Width)

	if t.Focused {
		boxStyle = boxStyle.BorderForeground(secondaryColor)
	} else {
		boxStyle = boxStyle.BorderForeground(surface1)
	}

	var content string
	if t.Value == "" && !t.Focused {
		content = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(t.Placeholder)
	} else {
		if t.Focused {
			before := t.Value[:t.CursorPos]
			cursor := "󰍟"
			after := ""
			if t.CursorPos < len(t.Value) {
				after = t.Value[t.CursorPos:]
			}

			content = before +
				lipgloss.NewStyle().Foreground(primaryColor).Background(surface0).Render(cursor) +
				after
		} else {
			content = t.Value
		}
	}

	b.WriteString(boxStyle.Render(content))

	return b.String()
}

type NumberInput struct {
	TextInput
	Min int
	Max int
}

func NewNumberInput(prompt, placeholder string, width, min, max int) NumberInput {
	return NumberInput{
		TextInput: NewTextInput(prompt, placeholder, width),
		Min:       min,
		Max:       max,
	}
}

func (n *NumberInput) InsertChar(ch rune) {
	if ch >= '0' && ch <= '9' {
		n.TextInput.InsertChar(ch)
	}
}
