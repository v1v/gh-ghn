package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/v72/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	noPrompt := flag.Bool("no-prompt", false, "Skip confirmation prompts")
	flag.Parse()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Please set the GITHUB_TOKEN environment variable.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	allNotifications := getUnreadNotifications(client, ctx)

	for _, notification := range allNotifications {
		subject := notification.GetSubject()
		if subject.GetType() == "PullRequest" {
			// Extract PR metadata
			repoFullName := strings.TrimPrefix(notification.GetRepository().GetFullName(), "repos/")
			prURL := subject.GetURL()
			parts := strings.Split(prURL, "/")
			owner := parts[4]
			repo := parts[5]
			prNumber := parts[len(parts)-1]
			prBrowserFriendlyURL := apiToWebURL(prURL)

			ownerRepo := strings.Split(repoFullName, "/")
			pr, _, err := client.PullRequests.Get(ctx, ownerRepo[0], ownerRepo[1], atoi(prNumber))
			if err != nil {
				log.Printf("Error fetching PR %s: %v\n", prURL, err)
				continue
			}

			prTitle := pr.GetTitle()
			prIsMerged := pr.GetMerged()
			prIsClosed := pr.GetState() == "closed"
			notificationIsRead := notification.GetUnread() == false

			if prIsMerged || prIsClosed {
				if notificationIsRead {
					continue
				}
				fmt.Printf("🟡 PR: %s, Title: \"%s\", is merged or closed and notification will be marked as read\n", prBrowserFriendlyURL, prTitle)

				threadID := getNotificationThreadID(allNotifications, owner, repo, prNumber)
				if threadID != "" {
					fmt.Printf("  🟡 \033[33mAbout to mark related GH Notification with threadID: \"%s\" as *READ*\033[0m\n", threadID)

					if promptProceed(*noPrompt) {
						markNotificationRead(client, ctx, threadID)
					}
				}
			} else {
				fmt.Printf("PR: %s, Title: \"%s\">, is unmerged and waiting for your review!\n", prBrowserFriendlyURL, pr.GetTitle())
			}

		}
	}
}

func getUnreadNotifications(client *github.Client, ctx context.Context) []*github.Notification {
	var allNotifications []*github.Notification
	opt := &github.NotificationListOptions{
		All:         false, // only include unread notifications
		ListOptions: github.ListOptions{PerPage: 50},
	}

	totalFetched := 0
	page := 1

	fmt.Printf("Fetching all GitHub notifications\n")
	for {
		opt.Page = page
		notifications, resp, err := client.Activity.ListNotifications(ctx, opt)
		if err != nil {
			log.Fatalf("Error fetching GH notifications: %v", err)
		}

		allNotifications = append(allNotifications, notifications...)
		totalFetched += len(notifications)

		fmt.Printf("\r📦 Page %d — Total fetched: %d notifications", page, totalFetched)

		if resp.NextPage == 0 {
			break
		}
		page++
	}

	fmt.Printf("\n✅ Done! Fetched %d notifications in total.\n", totalFetched)
	return allNotifications
}

func markNotificationRead(client *github.Client, ctx context.Context, threadID string) {
	_, err := client.Activity.MarkThreadRead(ctx, threadID)
	if err != nil {
		log.Printf("    ❌ Failed to mark thread %s as read: %v\n", threadID, err)
	} else {
		fmt.Printf("    🟢 Successfully marked thread as read\n")
	}
}

func getNotificationThreadID(notifications []*github.Notification, owner string, repo string, prNumber string) string {
	for _, n := range notifications {
		if strings.Contains(n.GetSubject().GetURL(), fmt.Sprintf("/repos/%s/%s/pulls/%s", owner, repo, prNumber)) {
			return n.GetID()
		}
	}

	return ""
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func promptProceed(skip bool) bool {
	if skip {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Continue: ? (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "", "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Invalid input, please enter 'y' or 'n', or just press Enter for yes.")
		}
	}
}

func apiToWebURL(apiURL string) string {
	// Convert API URL to GitHub Web URL
	url := strings.Replace(apiURL, "https://api.github.com/repos/", "https://github.com/", 1)
	return strings.Replace(url, "/pulls/", "/pull/", 1)
}
