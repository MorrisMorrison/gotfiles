package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	Dotfiles []string `json:"dotfiles"`
}

func loadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// processFile handles backing up a dotfile and ensuring the symlink exists.
// The isSync flag differentiates between init (false) and sync (true) for logging.
func processFile(file, homeDir, dotfilesRepoDir string, isSync bool) {
	sourcePath := filepath.Join(homeDir, file)
	destPath := filepath.Join(dotfilesRepoDir, file)

	// Check if the file exists in home.
	fi, err := os.Lstat(sourcePath)
	if err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			fmt.Printf("Skipping backup for %s as it is already a symlink.\n", file)
		} else {
			// Not a symlink: copy and remove original.
			if err := copyFile(sourcePath, destPath); err != nil {
				log.Printf("Error copying %s: %v", file, err)
			} else {
				if isSync {
					fmt.Printf("Updated %s in repository.\n", file)
				} else {
					fmt.Printf("Copied %s to repository.\n", file)
				}
			}
			if err := os.Remove(sourcePath); err != nil {
				log.Printf("Error removing original file %s: %v", file, err)
			}
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Error accessing %s: %v", file, err)
	} else {
		log.Printf("%s does not exist in home.", file)
	}

	// Ensure the symlink exists.
	if _, err := os.Lstat(sourcePath); os.IsNotExist(err) {
		if _, err := os.Stat(destPath); err == nil {
			if err := os.Symlink(destPath, sourcePath); err != nil {
				log.Printf("Error creating symlink for %s: %v", file, err)
			} else {
				fmt.Printf("Created symlink for %s.\n", file)
			}
		} else {
			log.Printf("No backup for %s found in repository.", file)
		}
	}
}

func runGitCommand(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

func initCmd(dotfiles []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}
	dotfilesRepoDir := filepath.Join(repoDir, "dotfiles")
	if err := os.MkdirAll(dotfilesRepoDir, 0755); err != nil {
		return err
	}

	for _, file := range dotfiles {
		processFile(file, homeDir, dotfilesRepoDir, false)
	}

	if err := runGitCommand(repoDir, "add", "."); err != nil {
		log.Printf("Error running git add: %v", err)
	}
	if err := runGitCommand(repoDir, "commit", "-m", "Update dotfiles backup"); err != nil {
		log.Printf("Error running git commit: %v", err)
	}
	if err := runGitCommand(repoDir, "push"); err != nil {
		log.Printf("Error running git push: %v", err)
	}
	return nil
}

func syncCmd(dotfiles []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}
	dotfilesRepoDir := filepath.Join(repoDir, "dotfiles")
	if _, err := os.Stat(dotfilesRepoDir); os.IsNotExist(err) {
		return fmt.Errorf("dotfiles repository directory does not exist. Run 'gotfiles init' first")
	}

	for _, file := range dotfiles {
		processFile(file, homeDir, dotfilesRepoDir, true)
	}

	if err := runGitCommand(repoDir, "add", "."); err != nil {
		log.Printf("Error running git add: %v", err)
	}
	if err := runGitCommand(repoDir, "commit", "-m", "Sync dotfiles changes"); err != nil {
		log.Printf("Error running git commit: %v", err)
	}
	if err := runGitCommand(repoDir, "push"); err != nil {
		log.Printf("Error running git push: %v", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gotfiles <init|sync>")
		os.Exit(1)
	}

	repoDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(repoDir, "config.json")
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config file (%s): %v", configPath, err)
	}

	switch os.Args[1] {
	case "init":
		if err := initCmd(cfg.Dotfiles); err != nil {
			log.Fatal(err)
		}
	case "sync":
		if err := syncCmd(cfg.Dotfiles); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
