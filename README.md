Here’s an improved **README.md** tailored for your real repository:

📦 **[https://github.com/Seann-Moser/hypr-u](https://github.com/Seann-Moser/hypr-u)**

I incorporated:

* Your repo name
* Your current config example
* Instructions for adding `hypr-u` to `exec-once` in Hyprland
* Clean layout suitable for GitHub
* Optional badges/header section

Feel free to copy/paste this directly into your GitHub repo.

---

# **Hypr-U — A File-Triggered Command Runner for Hyprland**

### *Automatically restart Hyprland components when their configs change.*

**Hypr-U** is a lightweight Go daemon that watches your Hyprland configuration files (or any files you choose) and automatically runs commands when they change.

Perfect for users who frequently edit configs and want instant reloads of:

* `waybar`
* `hypridle`
* `hyprpaper`
* custom scripts
* *anything* you want to auto-restart

Hypr-U is fast, minimal, self-reloading, and extremely configurable.

---

## ✨ Features

* 🚀 Automatically **reload Hyprland components** when configs change
* 🔄 **Multiple commands per file** (stop → start sequences)
* ⚙️ **Background or foreground** execution per command
* 🔍 Watches **any file or directory**
* 🧠 Watches and reloads **its own config file** automatically
* 💾 Creates a default `~/.config/hypr-u.yaml` if missing
* 📁 Supports `~` expansion in paths
* 🧵 Non-blocking, channel-based event system
* 📦 Fully self-contained (default config embedded with `go:embed`)

---

## 📦 Installation

Clone and build:

```sh
git clone https://github.com/Seann-Moser/hypr-u
cd hypr-u
go build -o hypr-u .
```

Install:

```sh
sudo cp hypr-u /usr/local/bin/
```

Or run directly:

```sh
./hypr-u
```

---

## 🏁 Autostart in Hyprland

Add this line to your Hyprland config:

```ini
exec-once = hypr-u
```

This ensures Hypr-U runs in the background and reloads your Hyprland components whenever configs change.

---

## 📄 Configuration

Hypr-U reads:

```
~/.config/hypr-u.yaml
```

If this file does **not** exist, it will automatically create it using the embedded default config.

### 🔁 Live Reloading

Hypr-U watches its own config file and automatically reloads whenever you edit it.
You never need to restart `hypr-u`.

---

## 🛠 YAML Format

### Top-Level Fields

```yaml
interval: 2s     # polling interval  
files:           # files to watch  
```

### Per-File Configuration

Each watched file has:

| Field                   | Description                       |
| ----------------------- | --------------------------------- |
| `path`                  | File to watch                     |
| `commands`              | List of commands to run on change |
| `commands[].path`       | Executable to run                 |
| `commands[].args`       | Arguments for command             |
| `commands[].background` | Run in background?                |

---

## 🔧 Example Configuration (your current default)

```yaml
# Default polling interval
interval: 2s

# Default files
files:
  - path: ~/.config/hypr/hypridle.conf
    commands:
      - path: killall
        args: ["hypridle"]
        background: false

      - path: /usr/bin/hypridle
        background: true

  - path: ~/.config/waybar/**
    commands:
      - path: killall
        args: ["waybar"]
        background: false

      - path: /usr/bin/waybar
        background: true
```

### What this does:

#### When `hypridle.conf` changes:

1. Runs `killall hypridle`
2. Starts `hypridle` in background

#### When `waybar/config.jsonc` changes:

1. Runs `killall waybar`
2. Starts `waybar` in background

---

## 🚀 Running Hypr-U

Run manually:

```sh
hypr-u
```

Or let Hyprland launch it:

```ini
exec-once = hypr-u
```

You'll see logs like:

```
Watching files...
Change detected in ~/.config/waybar/config.jsonc
Ran command killall
Started background command /usr/bin/waybar
```

---

