package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TextInput represents a text input field
type TextInput struct {
	Prompt      string
	Placeholder string
	Value       string
	CursorPos   int
	Width       int
	Focused     bool
}

// NewTextInput creates a new text input
func NewTextInput(prompt, placeholder string, width int) TextInput {
	return TextInput{
		Prompt:      prompt,
		Placeholder: placeholder,
		Width:       width,
		CursorPos:   0,
		Focused:     false,
	}
}

// InsertChar inserts a character at cursor position
func (t *TextInput) InsertChar(ch rune) {
	if t.CursorPos > len(t.Value) {
		t.CursorPos = len(t.Value)
	}
	
	before := t.Value[:t.CursorPos]
	after := t.Value[t.CursorPos:]
	t.Value = before + string(ch) + after
	t.CursorPos++
}

// DeleteChar deletes character before cursor
func (t *TextInput) DeleteChar() {
	if t.CursorPos > 0 {
		before := t.Value[:t.CursorPos-1]
		after := t.Value[t.CursorPos:]
		t.Value = before + after
		t.CursorPos--
	}
}

// MoveCursorLeft moves cursor left
func (t *TextInput) MoveCursorLeft() {
	if t.CursorPos > 0 {
		t.CursorPos--
	}
}

// MoveCursorRight moves cursor right
func (t *TextInput) MoveCursorRight() {
	if t.CursorPos < len(t.Value) {
		t.CursorPos++
	}
}

// MoveCursorStart moves cursor to start
func (t *TextInput) MoveCursorStart() {
	t.CursorPos = 0
}

// MoveCursorEnd moves cursor to end
func (t *TextInput) MoveCursorEnd() {
	t.CursorPos = len(t.Value)
}

// Clear clears the input
func (t *TextInput) Clear() {
	t.Value = ""
	t.CursorPos = 0
}

// Render renders the text input
func (t TextInput) Render() string {
	var b strings.Builder
	
	// Prompt
	promptStyle := lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true)
	b.WriteString(promptStyle.Render(t.Prompt))
	b.WriteString(" ")
	
	// Input box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(t.Width)
	
	if t.Focused {
		boxStyle = boxStyle.BorderForeground(secondaryColor)
	} else {
		boxStyle = boxStyle.BorderForeground(mutedColor)
	}
	
	// Value or placeholder
	var content string
	if t.Value == "" && !t.Focused {
		content = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(t.Placeholder)
	} else {
		// Show value with cursor
		if t.Focused {
			before := t.Value[:t.CursorPos]
			cursor := "█"
			after := ""
			if t.CursorPos < len(t.Value) {
				after = t.Value[t.CursorPos:]
			}
			
			content = before + 
				lipgloss.NewStyle().Foreground(secondaryColor).Render(cursor) + 
				after
		} else {
			content = t.Value
		}
	}
	
	b.WriteString(boxStyle.Render(content))
	
	return b.String()
}

// NumberInput represents a number input field
type NumberInput struct {
	TextInput
	Min int
	Max int
}

// NewNumberInput creates a new number input
func NewNumberInput(prompt, placeholder string, width, min, max int) NumberInput {
	return NumberInput{
		TextInput: NewTextInput(prompt, placeholder, width),
		Min:       min,
		Max:       max,
	}
}

// InsertChar only allows digits
func (n *NumberInput) InsertChar(ch rune) {
	if ch >= '0' && ch <= '9' {
		n.TextInput.InsertChar(ch)
	}
}
