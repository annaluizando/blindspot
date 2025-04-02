package cli

type Config struct {
	Difficulty int

	flagsChanged map[string]bool
}

// creates new CLI configuration with defaults
func NewConfig() *Config {
	return &Config{
		Difficulty:   0,
		flagsChanged: make(map[string]bool),
	}
}

// sets game difficulty level that is passed through "-d" in cli mode
func (c *Config) SetDifficulty(level int) {
	if level < 0 {
		level = 0
	} else if level > 2 {
		level = 2
	}
	c.Difficulty = level
}

// records whether a specific flag was changed on the command line
func (c *Config) SetFlagChanged(flag string, changed bool) {
	c.flagsChanged[flag] = changed
}

// checks if a flag was explicitly set on the command line
func (c *Config) WasFlagChanged(flag string) bool {
	return c.flagsChanged[flag]
}
