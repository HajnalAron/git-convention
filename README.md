# Git Convention CLI

Git Convention CLI is a command-line tool designed to help developers create standardized branch names and commit messages following conventional formats. This tool aims to improve consistency in Git workflows and make it easier to follow best practices.

## Features

- Generate conventional branch names
- Create structured commit messages
- Customizable branch and commit types via configuration
- Easy-to-use command-line interface

## Installation

To install Git Convention CLI, follow these steps:

1. Ensure you have Go installed on your system (version 1.16 or later).
2. Run the following command:

```
go install github.com/hajnalaron/git-convention-cli@latest
```

## Usage

Git Convention CLI provides the following commands:

- `branch`: Create a conventional branch name
- `commit`: Generate a structured commit message
- `config`: View or edit the configuration
- `help`: Display help information about any command

### Basic Usage

```
git-convention-cli [command] [flags]
```

### Available Commands

1. Create a branch:
   ```
   git-convention-cli branch
   ```

2. Create a commit message:
   ```
   git-convention-cli commit
   ```

3. View configuration:
   ```
   git-convention-cli config view
   ```

4. Get help:
   ```
   git-convention-cli help [command]
   ```

### Flags

- `--config`: Specify a custom configuration file path

Example:
```
git-convention commit --config /path/to/custom/config.json
```

## Configuration

Git Convention CLI uses a JSON configuration file to customize its behavior. By default, it looks for a configuration file in the following location:

```
~/.config/git-convention-cli/config.json
```

If no configuration file is found, a default one will be created at this location.

### Custom Configuration

You can specify a custom configuration file using the `--config` flag:

```
git-convention branch --config /path/to/custom/config.json
```

### Configuration Format

The configuration file is in JSON format and includes the following main sections:

- `default_branch_prefix`: The default prefix for branch names
- `default_commit_prefix`: The default prefix for commit messages
- `emojis_enabled`: Boolean to enable or disable emoji usage
- `branch_types`: Array of objects defining branch types
- `commit_types`: Array of objects defining commit types

Each branch type and commit type object contains:
- `type`: The short identifier for the type
- `description`: A brief description of the type
- `emoji`: The associated emoji (if emojis are enabled)
- `prefix`: The prefix used in branch names (for branch types only)

Example configuration:

```json
{
  "default_branch_prefix": "feature",
  "default_commit_prefix": "feat",
  "emojis_enabled": true,
  "branch_types": [
    {
      "type": "feature",
      "description": "A new feature",
      "emoji": "✨",
      "prefix": "feat"
    },
    {
      "type": "bugfix",
      "description": "A bug fix",
      "emoji": "🐛",
      "prefix": "fix"
    }
  ],
  "commit_types": [
    {
      "type": "feat",
      "description": "A new feature",
      "emoji": "✨"
    },
    {
      "type": "fix",
      "description": "A bug fix",
      "emoji": "🐛"
    }
  ]
}
```

This configuration allows for extensive customization of branch and commit types, including descriptions and emojis for each type. You can add, remove, or modify types to suit your project's needs.

## Contributing

Contributions to Git Convention CLI are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any problems or have any questions, please open an issue on the GitHub repository.

---

Happy conventional Git-ing!
