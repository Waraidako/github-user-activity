package main

import (
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

type CompressedEvent struct {
	Count int
	Event Event
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	if len(os.Args) > 2 {
		return fmt.Errorf("too many arguments")
	}

	events, err := fetchEvents(os.Args[1])
	if err != nil {
		return fmt.Errorf("fetch failure: %s", err)
	}
	if len(events) == 0 {
		fmt.Printf("No recent activity found\n")
		return nil
	}

	compressedEvents := compressEvents(events)
	lines := prepareOutputLines(compressedEvents)

	fmt.Printf("\nRecent events for %s:\n\n", os.Args[1])
	for _, line := range lines {
		fmt.Printf("  - %s\n", line)
	}
	fmt.Println()

	return nil
}

func printHelp() {
	lines := []string{
		"",
		"Github User Activity CLI",
		"",
		"Usage: github-activity <username>",
		"Outputs the user's activity if user exists",
		"",
	}
	fmt.Println(strings.Join(lines, "\n"))
}

func fetchEvents(username string) (events []Event, err error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	endpoint := fmt.Sprintf("https://api.github.com/users/%s/events", url.PathEscape(username))
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "github-user-activity-cli")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK response status: %s", response.Status)
	}

	if err := json.UnmarshalRead(response.Body, &events); err != nil {
		return nil, fmt.Errorf("unable to parse response: %s", err)
	}
	return events, nil

}

func compressEvents(events []Event) []CompressedEvent {
	compressedEvents := []CompressedEvent{
		{
			Count: 1,
			Event: events[0],
		},
	}

	for i := 1; i < len(events); i++ {
		lastCompressedEvent := &compressedEvents[len(compressedEvents)-1]
		if events[i].Repo.Name == lastCompressedEvent.Event.Repo.Name && events[i].Type == lastCompressedEvent.Event.Type {
			lastCompressedEvent.Count++
		} else {
			compressedEvents = append(compressedEvents, CompressedEvent{Count: 1, Event: events[i]})
		}
	}
	return compressedEvents
}

func prepareOutputLines(events []CompressedEvent) []string {
	lines := make([]string, 0)
	for _, event := range events {
		lines = append(lines, constructLine(event))
	}
	return lines
}

func constructLine(event CompressedEvent) string {
	switch event.Event.Type {
	case "PushEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Pushed to %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Pushed to %s %d times", event.Event.Repo.Name, event.Count)

	case "CreateEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Created a new branch, tag, or repo in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Created a branch, tag, or repo %d times in %s", event.Count, event.Event.Repo.Name)

	case "DeleteEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Deleted a tag or a branch in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Deleted %d tags or branches in %s", event.Count, event.Event.Repo.Name)

	case "IssuesEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Interacted with an issue in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Interacted with an issue %d times in %s", event.Count, event.Event.Repo.Name)

	case "IssueCommentEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Commented on an issue or PR in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Commented on an issue or PR %d times in %s", event.Count, event.Event.Repo.Name)

	case "PullRequestEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Modified a PR in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Modified PRs %d times in %s", event.Count, event.Event.Repo.Name)

	case "PullRequestReviewEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Performed a PR review in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Performed %d PR reviews in %s", event.Count, event.Event.Repo.Name)

	case "PullRequestReviewCommentEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Commented on a PR review in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Commented on %d PR reviews in %s", event.Count, event.Event.Repo.Name)

	case "CommitCommentEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Commented on a commit in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Commented on commits %d times in %s", event.Count, event.Event.Repo.Name)

	case "ForkEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Forked %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Forked %s %d times", event.Event.Repo.Name, event.Count)

	case "WatchEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Starred %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Starred %s %d times", event.Event.Repo.Name, event.Count)

	case "ReleaseEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Published a release in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Published %d releases in %s", event.Count, event.Event.Repo.Name)

	case "DiscussionEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Created a discussion in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Created %d discussions in %s", event.Count, event.Event.Repo.Name)

	case "GollumEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Edited wiki in %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Edited wiki %d times in %s", event.Count, event.Event.Repo.Name)

	case "MemberEvent":
		if event.Count == 1 {
			return fmt.Sprintf("Added a collaborator to %s", event.Event.Repo.Name)
		}
		return fmt.Sprintf("Added %d collaborators to %s", event.Count, event.Event.Repo.Name)

	case "PublicEvent":
		return fmt.Sprintf("Made %s public", event.Event.Repo.Name)

	default:
		return fmt.Sprintf("Performed an unknown event in %s", event.Event.Repo.Name)
	}
}
