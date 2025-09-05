package ui

import "github.com/charmbracelet/bubbles/key"

// Common key bindings used across multiple UI components
var (
	// Navigation keys
	Up = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	)
	Down = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	)
	ScrollUp = key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "scroll up"),
	)
	ScrollDown = key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "scroll down"),
	)
	Select = key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter", "select"),
	)
	Back = key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	)
	Next = key.NewBinding(
		key.WithKeys("n", "enter"),
		key.WithHelp("n", "next"),
	)
	Help = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	)
	Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	)
	ShowHint = key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "show hint"),
	)
)

// MenuKeyMap defines key bindings for menu components
type MenuKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Select     key.Binding
	Back       key.Binding
	Help       key.Binding
	Quit       key.Binding
}

// NewMenuKeyMap creates a new menu key map with common bindings
func NewMenuKeyMap() MenuKeyMap {
	return MenuKeyMap{
		Up:         Up,
		Down:       Down,
		ScrollUp:   ScrollUp,
		ScrollDown: ScrollDown,
		Select:     Select,
		Back:       Back,
		Help:       Help,
		Quit:       Quit,
	}
}

// key bindings for challenge view
type ChallengeKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Select     key.Binding
	Back       key.Binding
	Help       key.Binding
	Quit       key.Binding
	ShowHint   key.Binding
	Next       key.Binding
}

func NewChallengeKeyMap() ChallengeKeyMap {
	return ChallengeKeyMap{
		Up:         Up,
		Down:       Down,
		ScrollUp:   ScrollUp,
		ScrollDown: ScrollDown,
		Select:     Select,
		Back:       Back,
		Help:       Help,
		Quit:       Quit,
		ShowHint:   ShowHint,
		Next:       Next,
	}
}

func NewExplanationKeyMap() ExplanationKeyMap {
	return ExplanationKeyMap{
		ScrollUp:   ScrollUp,
		ScrollDown: ScrollDown,
		Next:       Next,
		Back:       Back,
		Help:       Help,
		Quit:       Quit,
	}
}
