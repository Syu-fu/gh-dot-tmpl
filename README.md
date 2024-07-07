# gh-dot-tmpl

[![lint](https://github.com/Syu-fu/gh-dot-tmpl/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/Syu-fu/gh-dot-tmpl/actions/workflows/lint.yml)
[![check license](https://github.com/Syu-fu/gh-dot-tmpl/actions/workflows/license-check.yml/badge.svg?branch=main)](https://github.com/Syu-fu/gh-dot-tmpl/actions/workflows/license-check.yml)
[![test](https://github.com/Syu-fu/gh-dot-tmpl/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/Syu-fu/gh-dot-tmpl/actions/workflows/test.yml)
[![Go Coverage](https://github.com/Syu-fu/gh-dot-tmpl/wiki/coverage.svg)](https://raw.githack.com/wiki/Syu-fu/gh-dot-tmpl/coverage.html)

`gh-dot-tmpl` is a GitHub CLI extension that generates the contents of the `.github` folder from templates.  
It helps in automating and standardizing the setup of GitHub repository configurations.

## Dependencies

`gh-dot-tmpl` requires the following dependencies:

- `gh` (GitHub CLI) version 2.51.0 or higher
- `git` version 2.45.2 or higher

Ensure that these dependencies are installed and properly configured before using `gh-dot-tmpl`.

## Installation

```shell
gh extension install Syu-fu/gh-dot-tmpl
```

Upgrade:

```shell
gh extension upgrade dot-tmpl
```

## Usage

### Running the Command

1. **Generate files from templates:**

To generate the contents of the `.github` folder from specified templates, use the following command:

```sh
gh dot-tmpl [TEMPLATE_NAME1] [TEMPLATE_NAME2] ...
```

Replace [TEMPLATE_NAME1], [TEMPLATE_NAME2], etc., with the names of the templates you want to use.

### Command Flags

| Flag          | Description                 |
| ------------- | --------------------------- |
| -h, --help    | Display help information    |
| -v, --version | Display version information |

### Configuration

The location for the configuration file is `$XDG_CONFIG_HOME/gh-dot-tmpl/config.yaml`.

#### Configuration File Example

Below is an example of a configuration file (config.yaml):

```yaml
templates:
  issue:
    template_file: ~/.config/gh-dot-tmpl/template/issue.md
    output_file: .github/ISSUE_TEMPLATE.md
  pr:
    template_file: ~/.config/gh-dot-tmpl/template/pullrequest.md
    output_file: .github/PULL_REQUEST_TEMPLATE.md
```

| Key           | Description                                                                           |
| ------------- | ------------------------------------------------------------------------------------- |
| templates     | A mapping of template names to their respective template files and output file names. |
| template_file | The name of the template file to use.                                                 |
| output_file   | The name of the file to generate.                                                     |

### Templates

Template files should be placed under `$XDG_CONFIG_HOME/gh-dot-tmpl/template/`.

Template Replacements
The following placeholders can be used in template files and will be replaced accordingly:

| Placeholder     | Description                        |
| --------------- | ---------------------------------- |
| {{.Username}}   | Replaced with the GitHub username. |
| {{.Repository}} | Replaced with the repository name. |

For example, a template file (issue.md) might look like this:

```md
---
name: Bug report
about: Create a report to help us improve
title: ""
labels: ""
assignees: ""
---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:

1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Information:**

- OS:
- {{.Reponame}} Version:

**Additional context**
Add any other context about the problem here.
```

## Contributing

We welcome contributions to gh-dot-tmpl! Please see the [CONTRIBUTING.md](https://github.com/Syu-fu/gh-dot-tmpl/blob/main/.github/CONTRIBUTING.md) file for guidelines on how to contribute to this project.

## License

Distributed under the [MIT License](https://github.com/Syu-fu/gh-dot-tmpl/blob/main/LICENSE).
