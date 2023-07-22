# Install

```shell
go install github.com/dushaoshuai/issues-blog@latest
```

# What is issues-blog?

Write blogs using GitHub issues. Sync blogs to and from GitHub issues.

What it can do:

* pull remote issue(s) to the local
* push a local blog to GitHub issues

What it can not do:

* create a GitHub issue
* delete a GitHub issue
* many other things ...

# Why issues-blog?

I write blogs (notes) using GitHub issues. After a few writings,
I found that the issues web editor is not very convenient and efficient for writing long articles.
I started writing locally and copying the contents to the web editor.
Finally, I decided to write a CLI to make things simpler.

# Usage

The default config file is `$HOME/.config/.isblog.yaml`. An example config file:

```yaml
owner: dushaoshuai
repo: dushaoshuai.github.io
token: <MY-TOKEN>
```

Usage:

```bash
$ isblog --help 
Write blogs using Github issues.
Sync blogs to and from Github issues.

Usage:
  isblog [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  pull        Pull the remote Github issue(s) to the local blog(s)
  push        Push the local blog to the remote Github issue

Flags:
      --config string   config file (default $HOME/.config/.isblog.yaml)
  -h, --help            help for isblog
  -v, --version         version for isblog

Use "isblog [command] --help" for more information about a command.
```
