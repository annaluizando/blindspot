package cli

type Config struct {
	Difficulty int
	Category   string
	flagsSet   map[string]bool
}

func NewConfig() *Config {
	return &Config{
		Difficulty: 0,
		Category:   "",
		flagsSet:   make(map[string]bool),
	}
}

func (c *Config) SetDifficulty(level int) {
	if level < 0 {
		level = 0
	} else if level > 2 {
		level = 2
	}
	c.Difficulty = level
}

func (c *Config) SetCategory(category string) {
	c.Category = category
}

func (c *Config) SetFlag(flag string, changed bool) {
	c.flagsSet[flag] = changed
}

func (c *Config) WasFlagSet(flag string) bool {
	return c.flagsSet[flag]
}
