package filemonitor

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

import _ "embed"

//go:embed config/default_config.yaml
var defaultConfig []byte
var ConfigPath = "~/.config/hypr-u.yaml"

// LoadConfig loads YAML into Config
func LoadConfig(path string) (*Config, error) {
	expanded, err := expandHome(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(expanded)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Interval == 0 {
		cfg.Interval = 5 * time.Second
	}

	return &cfg, nil
}

func NewFileWatcher(cfg *Config) (*FileWatcher, error) {
	cfgPath, _ := expandHome(ConfigPath)

	// Stat config file
	st, err := os.Stat(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("cannot stat config file: %w", err)
	}

	fw := &FileWatcher{
		Files:      make(map[string]*fileState),
		FileCfgs:   make(map[string]FileConfig),
		Interval:   cfg.Interval,
		Events:     make(chan FileChangeEvent, 10),
		ConfigFile: cfgPath,
		ConfigState: fileState{
			FullPath: cfgPath,
			ModTime:  st.ModTime(),
			Size:     st.Size(),
		},
	}

	// Load file entries from config
	for _, f := range cfg.Files {
		p, _ := expandHome(f.Path)
		st, err := os.Stat(p)
		if err != nil {
			continue
		}
		fw.Files[p] = &fileState{
			FullPath: p,
			ModTime:  st.ModTime(),
			Size:     st.Size(),
		}
		fw.FileCfgs[p] = f
	}

	return fw, nil
}

// Start polling files
func (fw *FileWatcher) Start() {
	ticker := time.NewTicker(fw.Interval)
	defer ticker.Stop()

	for range ticker.C {
		fw.checkFiles()
	}
}

func (fw *FileWatcher) checkFiles() {
	// Check normal files
	for path, state := range fw.Files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.ModTime() != state.ModTime || info.Size() != state.Size {
			state.ModTime = info.ModTime()
			state.Size = info.Size()

			fw.Events <- FileChangeEvent{
				Path:      path,
				Config:    fw.FileCfgs[path],
				Timestamp: time.Now(),
			}
		}
	}

	// Check config file
	cfgInfo, err := os.Stat(fw.ConfigFile)
	if err == nil {
		if cfgInfo.ModTime() != fw.ConfigState.ModTime || cfgInfo.Size() != fw.ConfigState.Size {
			fw.ConfigState.ModTime = cfgInfo.ModTime()
			fw.ConfigState.Size = cfgInfo.Size()

			// Send special event
			fw.Events <- FileChangeEvent{
				Path:      fw.ConfigFile,
				Timestamp: time.Now(),
				Config:    FileConfig{}, // no command
			}
		}
	}
}

// Expand home directory path ("~")
func expandHome(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, path[1:]), nil
}
func createDefaultConfig(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, defaultConfig, 0o644)
}

func LoadDefaultConfig() (*Config, error) {
	configPath, err := expandHome("~/.config/hypr-u.yaml")
	if err != nil {
		return nil, err
	}

	// If config does not exist, create it
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return nil, err
		}
	}

	return LoadConfig(configPath)
}

func (fw *FileWatcher) StartWorker() {
	go func() {
		for ev := range fw.Events {

			// CONFIG RELOAD TRIGGER
			if ev.Path == fw.ConfigFile {
				log.Println("Config file changed — reloading...")
				if err := fw.ReloadConfig(); err != nil {
					log.Printf("Failed to reload config: %v", err)
				}
				log.Printf("done")
				continue
			}

			// NORMAL FILE ACTION
			cfg := ev.Config
			if len(cfg.Cmd) == 0 {
				log.Printf("Change detected in %s (no commands)", ev.Path)
				continue
			}

			for _, c := range cfg.Cmd {

				cmd := exec.Command(c.Path, c.Args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if c.Background {
					cmd.Stdout = io.Discard
					cmd.Stderr = io.Discard
					// Non-blocking: run command in background
					if err := cmd.Start(); err != nil {
						log.Printf("Failed to start background command %s: %v", c.Path, err)
					} else {
						log.Printf("Started background command %s", c.Path)
					}
					continue
				}

				// Blocking execution (foreground)
				if err := cmd.Run(); err != nil {
					log.Printf("Command execution failed: %v", err)
				} else {
					log.Printf("Ran command %s", c.Path)
				}
			}

		}
	}()
}

func (fw *FileWatcher) ReloadConfig() error {
	cfg, err := LoadDefaultConfig()
	if err != nil {
		return err
	}

	// Clear & repopulate
	fw.Files = make(map[string]*fileState)
	fw.FileCfgs = make(map[string]FileConfig)

	for _, f := range cfg.Files {
		p, _ := expandHome(f.Path)
		st, err := os.Stat(p)
		if err != nil {
			continue
		}

		fw.Files[p] = &fileState{
			FullPath: p,
			ModTime:  st.ModTime(),
			Size:     st.Size(),
		}
		fw.FileCfgs[p] = f
	}

	fw.Interval = cfg.Interval
	return nil
}
