# gotfiles

Gotfiles is a simple Go-based tool for synchronizing your dotfiles and configuration files using Git. It helps you manage, backup, and deploy your dotfiles by creating a local repository, copying your current dotfiles, and establishing symlinks in your home directory.

## Features

- **Initial Backup:** Backs up your existing dotfiles from your home directory into a repository.
- **Symlink Creation:** Automatically removes the original dotfiles and creates symlinks that point to the repository.
- **Sync Changes:** Easily sync new or updated dotfiles to the repository.
- **Configurable:** Specify which dotfiles to track using a JSON configuration file.

## Requirements

- [Go](https://golang.org/) (version 1.13+ recommended)
- [Git](https://git-scm.com/)
- A Unix-like operating system for symlink support

## Installation

1. **Clone the Repository:**

   ```bash
   git clone <your-repo-url> gotfiles
   cd gotfiles
   ```

2. **Build the Tool:**

   ```bash
   go build -o gotfiles main.go
   ```

## Configuration

Create a `config.json` file in the root of the repository to specify the dotfiles you want to track. For example:

```json
{
  "dotfiles": [
    ".bashrc",
    ".vimrc",
    ".gitconfig",
    ".tmux.conf"
  ]
}
```

Modify this list to suit your personal setup.

## Usage

Gotfiles supports two main commands:

### `init`

Use this command on a system where you already have dotfiles:

- **Backup:** Copies your dotfiles from your home directory into the repository.
- **Replace:** Removes the original files and creates symlinks in your home directory that point to the repository copies.
- **Git Operations:** Stages, commits, and pushes the changes to your remote Git repository.

Run the command:

```bash
./gotfiles init
```

### `sync`

Use this command to sync any new or updated dotfiles or to set up a new system:

- **Update:** Checks for modifications in your home directory and updates the repository backup.
- **Symlink Check:** Ensures that symlinks exist for all the configured dotfiles.
- **Git Operations:** Stages, commits, and pushes the changes to your remote Git repository.

Run the command:

```bash
./gotfiles sync
```

## How It Works

1. **Configuration Loading:**  
   Gotfiles reads the list of dotfiles from `config.json`.

2. **Processing Files:**  
   For each dotfile, the tool:
   - Checks if the file exists in your home directory.
   - Copies it to the repository (if it's not already a symlink).
   - Removes the original file.
   - Creates a symlink pointing to the backup in the repository.

3. **Git Integration:**  
   After processing the files, Gotfiles automatically stages, commits, and pushes any changes to your Git repository.

## Contributing

Contributions are welcome! If you have any suggestions or improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```
