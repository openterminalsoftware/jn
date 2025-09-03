# jn
An open-source markdown previewer and note-taking command line interface/application written in Go from scratch.

[See development progress](TODO.md)

# Documentation
## Installation
Right now there are no pre-built [release](https://github.com/openterminalsoftware/jn/releases) binaries, so you'll need to build from source and install like so.

```bash
git clone https://github.com/openterminalsoftware/jn
cd jn
make install
```

Make sure you have a `$GOBIN` environment variable correctly configured to your local `bin`.
```bash
export GOBIN=/usr/local/bin/
```

> [!NOTE]
> Adjust `/usr/local/bin/` to your local `bin` directory.
> Also, to use `make` you need to install it.

## Configuration

Create a configuration file at `$HOME/.jn/config.json`, current supported fields:
* `vault`

```json
{
    "vault": "~/.jn/vault"
}

## Commands

All commands are run from your terminal.

### `new`

Create a new note. This command opens a simple text editor in your terminal where you can write your markdown content.

**Usage:**

```bash
jn new
```

After running the command, you will be prompted to enter your note content. To save and quit the editor, type `.exit` on a new line and press Enter.

You will then be prompted to enter a filename for your note. The `.md` extension will be added automatically if you don't include it.

**Example:**

```bash
$ jn new
Type your markdown ".exit" to save & quit
# My new note
This is a test note.
.exit

Enter a filename for your note: my-first-note
Entry saved at: /Users/username/.jn/vault/my-first-note.md
```

### `list`

List all notes in your vault.

**Usage:**

```bash
jn list
```

This command will display a numbered list of all markdown files in your vault directory (`~/.jn/vault` by default).

**Example:**

```bash
$ jn list
Notes in vault:
1. my-first-note.md
2. another-note.md
```

### `search`

Search for notes by filename or content. This command opens an interactive search prompt.

**Usage:**

```bash
jn search
```

As you type, the search results will update in real-time.

**Interactive Keys:**

*   **Enter:** Quit the search.
*   **Tab:** Preview the first search result.
*   **Backspace:** Delete the last character in your search query.

**Example:**

```bash
$ jn search
Search: first

--- Results ---
  1. my-first-note.md
    # My new note
```

### `delete`

Delete a note from your vault.

**Usage:**

```bash
jn delete [filename]
```

You can specify the filename with or without the `.md` extension.

**Example:**

```bash
$ jn delete my-first-note
Deleted /Users/dennis/.jn/vault/my-first-note.md
```

### `preview`

Preview a markdown file in your terminal with syntax highlighting.

**Usage:**

```bash
jn preview [path/to/file.md]
```

**Example:**

```bash
$ jn preview ~/.jn/vault/another-note.md
```
This will print the content of `another-note.md` to the terminal with markdown formatting.

### `version`

Show the version of the `jn` application.

**Usage:**

```bash
jn version
```

**Example:**

```bash
$ jn version
You are using jn 1.0.0
Created by OpenTerminalSoftware (OTS)
https://github.com/openterminalsoftware/jn
```
