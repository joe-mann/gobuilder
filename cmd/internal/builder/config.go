package builder

type Config struct {
	rootDir      string
	gopath       string
	excludedDirs []string
}

func NewConfig(rootDir, gopath string) *Config {
	return &Config{
		rootDir:      rootDir,
		gopath:       gopath,
		excludedDirs: []string{".git", "vendor"},
	}
}

func (c *Config) SetExcludedDirs(dirs []string) {
	c.excludedDirs = dirs
}
