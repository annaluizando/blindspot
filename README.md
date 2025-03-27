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

# 🛡️ Secure Code Game !

A CLI-based interactive game designed to train yourself to identify insecure coding practices based on the OWASP Top 10, in a practical and fun way!

Feel free to play, analyze, modify, and contribute. New ideas and collaborations are always welcome :)

## 🔧 Stack

- **Golang** - Core language
- [**BubbleTea**](https://github.com/charmbracelet/bubbletea) - Terminal UI framework

## ✨ Features

The following features can be accessed in the "Settings" section, allowing players to customize their learning experience:

### 🔍 Vulnerability Name Display Toggle

It's possible to change visibility of vulnerability names to adjust if those names should appear or not in the top of challenges.
Hiding names creates a more challenging experience where you must identify vulnerabilities completely on your own.

### 🔄 Challenge Order

Two playing modes are available:

- **Random by Difficulty** - A more "advanced" way of playing, where vulnerabilities appear grouped by difficulty level, progressing from beginner to advanced. If you want the "hardest" way of playing, combine this with Vulnerability names: Hide.
- **Category Order** - A more directed way to train your eye for specific vulnerability category, as vulnerabilities appear grouped by their category.

## 📝 Progress Tracking

Track your learning journey through each vulnerability category wit completion percentages for each category that can be seen in "Categories" and/or "Progress".

## ⌨️ Controls

- Use arrow keys to navigate
- Press `Enter` to select
- Press `n` or `Enter` to go to next challenge
- Press `q` to quit at any time

## To-Do

### Completed

- [x] Implement syntax highlighting
- [x] Add more challenges and vulnerabilities
- [x] Add explanations for each vulnerability category
- [x] Fix "Progress" tab functionality
- [x] Implement both random and category-based game modes
- [x] Restore "Settings" section with customization options
- [x] Add ability to skip challenges
- [x] Include vulnerability explanations in category sections
- [x] Fix some texts line break (menu screen)

### In Progress

- [ ] Add Section 2 for practical exercises
- [ ] Implement adaptive text color based on terminal theme
- [ ] Fix and improve help text clarity
- [ ] Complete manual review of all challenges
- [ ] Correct vuln explanation to go back to its origin and not main menu when user presses back key
- [ ] Add Congratulation screen when finishing all challenges

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
