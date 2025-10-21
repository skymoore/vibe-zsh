---
title: "Usage"
weight: 3
---


## Basic Usage

Using vibe is simple:

1. **Type** a natural language description in your terminal
2. **Press** `Ctrl+G`
3. **Review** the generated command (with explanations)
4. **Press** `Enter` to execute, or edit first

## Query History

Vibe automatically saves your queries and generated commands. Access your history in two ways:

### Using the Keybinding

Press `Ctrl+X` then `H` to open the interactive history menu.

### Using the Command

Type `vh` and press Enter to open the interactive history menu.

### Quick Regenerate

Press `Ctrl+X` then `G` to instantly regenerate a new command from your most recent query without opening the menu. This is useful when you want to quickly try a different variation of your last command.

### History Menu Features

- **Navigate** with arrow keys (or `j`/`k`)
- **Search** by pressing `/` to filter entries
- **Select** by pressing `Enter` - inserts the generated command into your buffer
- **Regenerate** by pressing `G` - generates a new command from the original query
- **Edit Query** by pressing `V` - puts the original query in your buffer for editing
- **Jump to Start** by pressing `A` or `Home`
- **Jump to End** by pressing `E` or `End`
- **Cancel** by pressing `Esc` or `Q`

**Use Cases:**
- **Enter**: Use the exact command that was generated before
- **G**: Get a fresh variation of the same query (might produce different results)
- **V**: Modify the original query before generating a new command

### History Management

```bash
# View history in plain text
vibe-zsh history list

# Clear all history
vibe-zsh history clear
```

**Note:** Don't run `vibe-zsh history` directly - use `vh` or `Ctrl+X H` for the interactive menu.

## Examples

### File Operations

```bash
show me all hidden files including their sizes
# Press Ctrl+G
# → ls -lah
```

```bash
find all python files modified today
# Press Ctrl+G
# → find . -name "*.py" -mtime 0
```

```bash
compress all log files into an archive
# Press Ctrl+G
# → tar -czf logs.tar.gz *.log
```

### Docker

```bash
show logs of nginx container and follow them
# Press Ctrl+G
# → docker logs -f nginx
```

```bash
list all docker containers including stopped ones
# Press Ctrl+G
# → docker ps -a
```

```bash
remove all stopped containers
# Press Ctrl+G
# → docker container prune
```

### Git

```bash
show me commits from last week
# Press Ctrl+G
# → git log --since="1 week ago"
```

```bash
undo the last commit but keep the changes
# Press Ctrl+G
# → git reset --soft HEAD~1
```

```bash
show me what changed in the last commit
# Press Ctrl+G
# → git show HEAD
```

### System Administration

```bash
show disk usage of all mounted filesystems
# Press Ctrl+G
# → df -h
```

```bash
find processes listening on port 8080
# Press Ctrl+G
# → lsof -i :8080
```

```bash
show memory usage sorted by process
# Press Ctrl+G
# → ps aux --sort=-%mem | head
```

### Network

```bash
test if port 443 is open on example.com
# Press Ctrl+G
# → nc -zv example.com 443
```

```bash
show my public ip address
# Press Ctrl+G
# → curl ifconfig.me
```

```bash
download a file and show progress
# Press Ctrl+G
# → curl -# -O <url>
```

## Advanced Features

### Tab Completion

Vibe provides tab completion for common queries:

```bash
vibe <TAB><TAB>
# Shows common query suggestions
```

### Interactive Mode

When `VIBE_INTERACTIVE=true`, vibe will ask for confirmation before inserting commands:

```bash
export VIBE_INTERACTIVE=true

# Type: delete all log files
# Press Ctrl+G

---
rm *.log
---
Execute this command? [Y/n]
```

### Command Explanations

By default, vibe shows explanations for each part of the command:

```bash
docker ps -a
# docker: Docker command-line tool
# ps: List containers
# -a: Show all containers (not just running)
```

To hide explanations:

```bash
export VIBE_SHOW_EXPLANATION=false
```

### Safety Warnings

Vibe will warn you about potentially dangerous commands:

```bash
# Type: delete everything in current directory
# Press Ctrl+G

⚠️  WARNING: This command may delete files
rm -rf *
```

To disable warnings:

```bash
export VIBE_SHOW_WARNINGS=false
```

## Tips and Best Practices

### Be Specific

More specific queries generate better commands:

❌ **Vague:** "list files"
✅ **Specific:** "list all files including hidden ones with human-readable sizes"

### Use Natural Language

Write how you would explain the task to someone:

✅ "show me all processes using more than 100MB of memory"
✅ "find all files larger than 1GB modified in the last hour"

### Review Before Executing

Always review the generated command before pressing Enter. Vibe is a tool to help you, not to blindly execute commands.

### Edit Generated Commands

The generated command is editable. Use arrow keys to modify it before executing:

```bash
# Generated: git log --since="1 week ago"
# Edit to:    git log --since="2 weeks ago" --author="yourname"
```

### Combine with Other Tools

Vibe-generated commands work with pipes and other shell features:

```bash
# Generate: docker ps
# Edit to:  docker ps | grep nginx
```

## Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+G` | Generate command from natural language |
| `Ctrl+X H` | Open interactive history menu |
| `Ctrl+X G` | Regenerate most recent command |

### Customizing Keybindings

To change the command generation keybinding, add this to your `.zshrc`:

```bash
bindkey '^X' vibe  # Use Ctrl+X instead of Ctrl+G
```

To change the history keybindings:

```bash
export VIBE_HISTORY_KEY="^R"      # Use Ctrl+R for history menu
export VIBE_REGENERATE_KEY="^[r"  # Use Alt+R for quick regenerate
```

**Note:** Avoid using `^H` (Ctrl+H) as it conflicts with Backspace.
