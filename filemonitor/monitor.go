package filemonitor

import "time"

type File struct {
	FileName    string `yaml:"fileName"`
	Path        string `yaml:"path"`
	ProgramPath string `yaml:"programPath"`
}

// Internal structure for watching

type FileConfig struct {
	Path string          `yaml:"path"`
	Cmd  []ProgramConfig `yaml:"commands"`
}

type ProgramConfig struct {
	Path       string   `yaml:"path"`
	Args       []string `yaml:"args"`
	Background bool     `yaml:"background"`
}

type Config struct {
	Files     []FileConfig  `yaml:"files"`
	Directory string        `yaml:"directory"`
	Interval  time.Duration `yaml:"interval"`
}

type FileChangeEvent struct {
	Path      string
	Config    FileConfig
	Timestamp time.Time
}

type fileState struct {
	FullPath string
	ModTime  time.Time
	Size     int64
}

type FileWatcher struct {
	Files       map[string]*fileState
	FileCfgs    map[string]FileConfig
	Interval    time.Duration
	Events      chan FileChangeEvent
	ConfigFile  string // full expanded path
	ConfigState fileState
}
