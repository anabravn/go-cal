package main

import (
  "context"
  "fmt"
  "log"
  "os"
  "time"

  "golang.org/x/oauth2/google"
  "google.golang.org/api/calendar/v3"
  "google.golang.org/api/option"
)

func GetUpcomingEvents(cal string, srv *calendar.Service) {
  t := time.Now().Format(time.RFC3339)
  events, err := srv.Events.List(cal).ShowDeleted(false).
    SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()

  if err != nil {
    log.Printf("Unable to retrieve upcoming events: %v", err)
    return
  }

  for _, item := range events.Items {
    fmt.Println(item.Summary)
  }

}


func GetCalendars(srv *calendar.Service) []string {

  calendars, err := srv.CalendarList.List().ShowDeleted(false). 
    MaxResults(10).Do()

  if err != nil {
    log.Fatalf("Unable to retrieve calendars: %v", err)
  }

  cal_id := make([]string, len(calendars.Items))

  for _, item := range calendars.Items {
    cal_id = append(cal_id, item.Id)
  }  

  return cal_id
}

func main() {
  ctx := context.Background()
  b, err := os.ReadFile("credentials.json")
  if err != nil {
    log.Fatalf("Unable to read client secret file: %v", err)
  }

  // If modifying these scopes, delete your previously saved token.json.
  config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
  if err != nil {
    log.Fatalf("Unable to parse client secret file to config: %v", err)
  }
  client := getClient(config)

  srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
  if err != nil {
    log.Fatalf("Unable to retrieve Calendar client: %v", err)
  }

  calendars := GetCalendars(srv)
  for _, id := range calendars {
      fmt.Println(id)
      GetUpcomingEvents(id, srv)
      fmt.Println()
  }
}

