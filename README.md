# 🔦 blindspot - a code security game!

A terminal interactive game designed to train yourself to identify insecure coding practices based on the OWASP Top 10, in a practical and fun way!

Feel free to play, analyze, modify, and contribute. New ideas and collaborations are always welcome :)

## 🪜 Get Started!

Here is step-by-step what you need to run this project:

### 📦 Installation

```
git clone https://github.com/annaluizando/blindspot.git
cd blindspot
make build
make install
```

### 🎮 Playing

Now, to play blindspot, you can simply run by it's name:

```
blindspot
```

#### 🚀 CLI Support

If you want to open blindspot directly in a random challenge, filtering by difficulty, you can run using flags.

##### 🎯 Available Flags:

| Flag | Short | Description | Values |
|------|-------|-------------|---------|
| `--difficulty` | `-d` | Set game difficulty level | `0` (Beginner), `1` (Intermediate), `2` (Advanced) |
| `--category` | `-c` | Set specific vulnerability category | Any category name from the game |

##### 🏆 Difficulty Levels:
- **`0` (Beginner)**: Perfect for newcomers to security concepts
- **`1` (Intermediate)**: For those with some security knowledge  
- **`2` (Advanced)**: For experienced security professionals

##### 💡 Examples:

**🎲 Playing random intermediate challenges:**
```bash
blindspot --difficulty=1
# or
blindspot -d 1
```

**🔒 Playing only injection challenges:**
```bash
blindspot --category="Injection"
# or
blindspot -c "Injection"
```

**⚡ Combining both flags:**
```bash
blindspot -d 1 -c "Cross-Site Scripting (XSS)"
```

**📋 See all available categories:**
```bash
blindspot --help
```

> ⚙️ How It Works:
When you use CLI flags, the game automatically:
1. **🎯 Difficulty Mode**: If `--difficulty` is set, switches to "Random by Difficulty" mode and filters challenges by the specified level
2. **📁 Category Mode**: If `--category` is set, jumps directly to the specified category and starts from the first challenge

## 🔧 Stack

- **Golang** - Core language
- [**BubbleTea**](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [**Cobra**](https://github.com/spf13/cobra) - for CLI support

## 📁 Project Structure

```
blindspot/
├── cmd/
│   └── game/
│       └── main.go                  # Entry point for the application
├── internal/
│   ├── challenges/                  # Package for all security challenges
│   │   ├── challenge.go             # Responsible for loading all challenges
│   │   └── vuln_explanation.go      # Responsible for loading all vulnerabilities explanation
|   ├── cli/
│   │   ├── config.go                # CLI configuration handlers
│   │   └── runner.go                # Program initialization and running
│   ├── ui/                          # UI components using Bubbletea
│   │   ├── styles.go                # Common styles/themes for the UI
│   │   ├── challenge_view.go        # View for displaying challenges
│   │   ├── vuln_explanation_view.go # View for displaying vulnerability explanation
│   │   ├── quiz_view.go             # Multiple choice component
│   │   └── menu.go                  # Navigation menu
│   ├── game/                        # Game logic
│   │   ├── state.go                 # Game state management
│   │   ├── progress.go              # User progress/statistics tracking
│   │   └── validator.go             # Code validation logic
│   └── utils/
│       └── highlight.go             # Syntax code highlighting
│       └── wrapText.go              # Text wrap according to terminal  width
├── assets/                          # Static assets
│   └── challenges.yaml              # Challenge definitions
│   └── vuln_explanations.yaml       # Vulnerabilities Explanations defitions
├── go.mod                           # Go module file
└── go.sum                           # Go dependencies
```

## ✨ Features

The following features can be accessed in the "Settings" section, which allows you to customize your learning experience:

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

## FAQ

### Q: I want to erase my progress and saved settings, how can I do this?

A: You can accomplish this by running the game > selecting "Settings" > clicking "Delete all progress data"

### Q: I want to write my own challenges, how can I do this?

A: You can accomplish this by changing your assets/challenges.yaml and filling information according
your questions. Each question needs to be in the following format:

```yaml
- category: string, all categories must have equal string to be in same category,
  description: string,
  challenges:

  - id: string,
    title: string, name of vulnerability,
    description: string, with brief description about problem,
    difficulty: 0 for Beginner, 1 for Intermediate and 2 for Advanced,
    code: |

    - your code pasted here -

    options:

    - "Option A"
    - "Option B"
    - "Option C"
      correctAnswer: string, equal to one of options,
      hint: string,
      lang: string, preferebly one of [supported languages in chroma lib](https://github.com/alecthomas/chroma?tab=readme-ov-file#supported-languages) for code highlighthing
```

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
- [x] Add Congratulation screen when finishing all challenges
- [x] Add progess data erasing in Settings
- [x] Correct vuln explanation to go back to its origin and not main menu when user presses back key
- [x] Add user errors count in each category for statistics
- [x] Add scroll when height is too big for terminal view
- [x] Fix and improve help text clarity
- [x] Complete manual review of challenges.yaml
- [x] Add cli support

### In Progress

- [ ] Resolve bug in game mode toggle
- [+-] Review code and improve

## 📜 License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.
