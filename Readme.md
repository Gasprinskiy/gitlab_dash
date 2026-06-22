# gitlab-dash

A terminal UI for monitoring GitLab projects. For each project listed in the
config, it shows how up-to-date the latest release tag is relative to the
target branch, and — optionally — whether there's a divergence between a test
branch and the target branch. All in a btop-like panel in your terminal.

## What it shows

- A list of projects defined in the config
- Tag freshness: whether the target branch has changes after the latest release
- (optional) Test branch divergence from the target branch, if a test branch is set
- Brief user information

## Requirements

- A GitLab Personal Access Token with API access (scope `read_api`)

## Configuration

The application reads a `config.yaml` file. Example:

```yaml
credentials:
  host: "https://gitlab.example.com"
  personal_token: {YOUR_PERSONAL_TOKEN}
projects_data:
  project_id_list: // projects ids you want to monitor
    - 1111
    - 2222
  test_branch_name: "test"
  ignore_test_branch_compare_list: // if some of the projects does not have test branch
    - 2222
```

The `test_branch_name` field is
optional: if set, the app shows a divergence flag between the test branch and
the target branch; if omitted, that check is skipped.


The application reads its config from the `.gitlab_dash` folder in your home
directory:

- **Linux / macOS** — `~/.gitlab_dash/config.yaml`
- **Windows** — `%USERPROFILE%\.gitlab_dash\config.yaml`

## How tag freshness is calculated

The latest tag of a project is compared against the target branch using
GitLab's compare API. If the target branch has commits after the tag, the
release is considered stale (the branch has accumulated unreleased changes).

The same applies to the test branch: its content is compared against the
target branch to show whether there's unmerged work.

## Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling and layout
- [Bubbles](https://github.com/charmbracelet/bubbles) — components (table)