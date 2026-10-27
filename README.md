# GitHub User Activity CLI

A lightweight command-line tool written in Go that fetches a GitHub user's recent activity and displays it in a readable format.

Built as a solution to the [GitHub User Activity project on roadmap.sh](https://roadmap.sh/projects/github-user-activity). No third-party libraries, API tokens, or configuration required.

## Features

- **Look up any GitHub username** and view their recent activity in the terminal.
- **Read human-friendly descriptions** of pushes, repository creation, issues, pull requests, reviews, stars, forks, releases, wiki edits, and other GitHub events.
- **Group consecutive matching events** by event type and repository, reducing repetitive output (for example, `Pushed to owner/repo 3 times`).
- **Handle errors** such as unknown users, unsuccessful API responses, network failures, and invalid argument counts.
- **Run without external dependencies**, using only Go's standard library and GitHub's REST API.

## Requirements

- **Go 1.27 or newer**, as specified in `go.mod` (the application uses `encoding/json/v2`).
- An internet connection to fetch activity from GitHub.

## Installation

Clone the repository and build the executable:

```sh
git clone https://github.com/Waraidako/github-user-activity.git
cd github-user-activity
go build .
```

This produces `github-activity` on Linux/macOS or `github-activity.exe` on Windows.

## Usage

Pass a GitHub username as the only argument:

```sh
./github-activity Waraidako
```

Example output (your results will vary with recent activity):

```text
Recent events for Waraidako:

  - Pushed to Waraidako/task-tracker-cli 3 times
  - Created a new branch, tag, or repo in Waraidako/task-tracker-cli
```

You can also run the project without building a binary first:

```sh
go run . Waraidako
```

Running the executable with no arguments displays usage instructions. Supplying more than one username or encountering a GitHub API failure returns an error.

## How it works

The application requests activity from the GitHub REST API:

```text
GET https://api.github.com/users/<username>/events
```

It decodes the JSON response, combines **adjacent events of the same type in the same repository**, and prints short descriptions. Counts refer to grouped API events, not necessarily the number of individual commits or distinct pull requests.

Activity is limited to what GitHub exposes through this endpoint; it is **not a complete account history**. Unauthenticated requests are subject to GitHub's API rate limits. The HTTP client uses a 10-second timeout.

## Project structure

```text
github-user-activity/
├── main.go       # HTTP requests, JSON parsing, grouping, and CLI output
├── go.mod
└── .gitignore
```

## Project reference

[GitHub User Activity — roadmap.sh](https://roadmap.sh/projects/github-user-activity)
