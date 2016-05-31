package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"

	"github.com/blang/whoshome"
)

type Calendar struct {
	Provider         whoshome.PresenceProvider
	CalendarID       string
	TickTime         time.Duration
	EventColorID     string
	ClientSecretFile string
	ClientTokenFile  string
	EventTitle       string
	PresentFunc      func(present []string) bool
	Threshold        int
	lastEventID      string
	lastEventStart   time.Time
	calSrv           *gcal.Service
	graceCounter     int
}

func (c *Calendar) InitAuth() error {
	b, err := ioutil.ReadFile(c.ClientSecretFile)
	if err != nil {
		return fmt.Errorf("Calendar: Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gcal.CalendarScope)
	if err != nil {
		return fmt.Errorf("Calendar: Unable to parse client secret file to config: %v", err)
	}
	tok, err := tokenFromFile(c.ClientTokenFile)
	if err == nil {
		return fmt.Errorf("Calendar: Token file already exists")
	}
	tok = getTokenFromWeb(config)
	saveToken(c.ClientTokenFile, tok)
	return nil
}

func (c *Calendar) Init() error {
	ctx := context.Background()
	b, err := ioutil.ReadFile(c.ClientSecretFile)
	if err != nil {
		return fmt.Errorf("Calendar: Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gcal.CalendarScope)
	if err != nil {
		return fmt.Errorf("Calendar: Unable to parse client secret file to config: %v", err)
	}
	tok, err := tokenFromFile(c.ClientTokenFile)
	if err != nil {
		return fmt.Errorf("Calendar: Can't read oauth2 token file: %v", err)
	}
	client := config.Client(ctx, tok)

	srv, err := gcal.New(client)
	if err != nil {
		return fmt.Errorf("Calendar: Unable to retrieve calendar Client %v", err)
	}

	c.calSrv = srv
	return nil
}

// Run runs in it's own GoRoutine
func (c *Calendar) Run() {
	ticker := time.Tick(c.TickTime)
	c.tick()
	for {
		select {
		case <-ticker:
			c.tick()
		}
	}
}

func (c *Calendar) tick() error {
	log.Printf("Tick")
	present, err := c.Provider.Present()
	if err != nil {
		return err
	}
	if c.PresentFunc(present) {
		c.graceCounter = 0
		if c.lastEventID != "" {
			log.Printf("Update event")
			c.updateEvent()
		} else {
			log.Printf("Create event")
			c.createEvent()
		}
	} else {
		log.Printf("Not Present")
		if c.graceCounter > c.Threshold {
			c.lastEventID = ""
		}
		if c.lastEventID != "" {
			c.graceCounter++
		}
	}
	return nil
}

func (c *Calendar) createEvent() error {
	now := time.Now()
	t := now.Format(time.RFC3339)
	t2 := now.Add(time.Minute).Format(time.RFC3339)
	event := &gcal.Event{
		Summary:     c.EventTitle,
		ColorId:     c.EventColorID,
		Location:    "",
		Description: "",
		Start: &gcal.EventDateTime{
			DateTime: t,
			TimeZone: "Europe/Berlin",
		},
		End: &gcal.EventDateTime{
			DateTime: t2,
			TimeZone: "Europe/Berlin",
		},
	}

	event, err := c.calSrv.Events.Insert(c.CalendarID, event).Do()
	if err != nil {
		return err
	}
	c.lastEventID = event.Id
	c.lastEventStart = now
	return nil
}
func (c *Calendar) updateEvent() error {
	now := time.Now()
	t := c.lastEventStart.Format(time.RFC3339)
	t2 := now.Format(time.RFC3339)
	event := &gcal.Event{
		Summary:     c.EventTitle,
		ColorId:     c.EventColorID,
		Location:    "",
		Description: "",
		Start: &gcal.EventDateTime{
			DateTime: t,
			TimeZone: "Europe/Berlin",
		},
		End: &gcal.EventDateTime{
			DateTime: t2,
			TimeZone: "Europe/Berlin",
		},
	}

	event, err := c.calSrv.Events.Update(c.CalendarID, c.lastEventID, event).Do()
	if err != nil {
		return err
	}
	c.lastEventID = event.Id
	return nil
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
