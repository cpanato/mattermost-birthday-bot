package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func main() {
	var flagConfigFile string
	flag.StringVar(&flagConfigFile, "config", "config-birthdaybot.json", "")
	flag.Parse()

	err := loadConfig(flagConfigFile)
	if err != nil {
		log.Fatalf("Error starting the job: %v", err)
	}

	calendarService, err := calendar.NewService(context.Background(), option.WithAPIKey(Config.GoogleCalendarAPIKey))
	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}

	currentTime := time.Now()
	timeMin := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
	timeMax := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, time.UTC)
	log.Println(timeMin.Format(time.RFC3339))
	log.Println(timeMax.Format(time.RFC3339))

	events, err := calendarService.Events.List(Config.GoogleCalendarID).ShowDeleted(false).
		SingleEvents(true).TimeMax(timeMax.Format(time.RFC3339)).TimeMin(timeMin.Format(time.RFC3339)).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	if len(events.Items) == 0 {
		log.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			attachment := []MMAttachment{}
			reg := regexp.MustCompile(`([^"]*).*-.*`)
			res := reg.ReplaceAllString(item.Summary, "${1}")
			msg := fmt.Sprintf("Happy Birthday **%v** :tada:", strings.TrimSpace(res))
			log.Println(msg)
			attach := MMAttachment{}
			attachment = append(attachment, *attach.AddField(MMField{Title: "Mattermost Birthdays", Value: msg}))

			payload := MMSlashResponse{
				Username:    "BirthdayBot",
				IconUrl:     "https://fcit.usf.edu/matrix/wp-content/uploads/2016/12/Robot-04-A.png",
				Attachments: attachment,
			}
			if Config.MMIncomingWebhook != "" {
				send(Config.MMIncomingWebhook, payload)
			}
		}
	}
}
