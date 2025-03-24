secure-code-game/
├── cmd/
│ └── security-game/
│ └── main.go # Entry point for the application
├── internal/
│ ├── challenges/ # Package for all security challenges
│ │ ├── challenge.go # Challenge interface definition
│ │ ├── injection/ # OWASP category: Injection, will be useful on section 2 (pratical exercises)
│ │ │ ├── sql.go # SQL injection challenges
│ │ │ └── command.go # Command injection challenges
│ │ ├── auth/ # OWASP category: Broken Authentication, will be useful on section 2 (pratical exercises)
│ │ │ └── password.go # Password-related challenges
│ │ └── ... # Other OWASP categories
│ ├── ui/ # UI components using Bubbletea
│ │ ├── styles.go # Common styles/themes for the UI
│ │ ├── challenge_view.go # View for displaying challenges
│ │ ├── editor_view.go # Code editor component, will be useful on section 2 (pratical exercises)
│ │ ├── quiz_view.go # Multiple choice component
│ │ └── menu.go # Navigation menu
│ ├── game/ # Game logic
│ │ ├── state.go # Game state management
│ │ ├── progress.go # User progress tracking
│ │ └── validator.go # Code validation logic
│ └── utils/ # Utility functions
│ └── highlight.go # Syntax highlighting
├── assets/ # Static assets
│ └── challenges.yaml # Challenge definitions
├── go.mod # Go module file
└── go.sum # Go dependencies

# Secure Code Game!

This is a cli game to promote a pratical and fun learning of secure code practices based on OWASP Top 10. :)
Feel free to analyze, play, modify, I'm accepting new ideas and collaboration!

## Stack

- Go
- BubbleTea (TUI)

## To-do:

- [x] Syntax highlight
- [x] Add more challenges/vulnerabilities
- [ ] Add section 2, for pratical exercises
- [ ] toggle text color from white to black depending on user's terminal color
- [x] add explanation from vulnerability in the beginning/ending of each category
- [x] "Settings" should be removed, probably
- [x] "Progress" tab not working, want to work and show how % of progress in each category user has
- [ ] the next exercise should be in the same level (beginner, intermediate).
- [x] revive "Settings" so user can choose if they want vulnerability name in top of file or not.
- [ ] user should be able to skip challenge using 'n' to next challenge
