package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func showEvents(srv *calendar.Service) {
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}

func main() {
	ctx := context.Background()
	client := NewClient()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	var cmdList = &cobra.Command{
		Use:   "show",
		Short: "show list",
		Run: func(cmd *cobra.Command, args []string) {
			showEvents(srv)
		},
	}
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdList)
	rootCmd.Execute()
}
