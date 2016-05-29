# Who's home Google Calendar

This application reads the linux arp file and creates google calendar events when a mac address is present.

Install it on your router and add your phones wifi mac address to track when you're home.

Calendar will update existing events for pretty looking time charts in your calender.

## Build

```
cd cmd/calendar

go build
```

## Usage

Create your client secrets using the [Google Console](https://console.developers.google.com/start/api?id=calendar) and download the credentials to `client-secret.json`.

First time use the `-init-token` flag to create an oauth2 token and save it to `-client-token-file`.
After that omit the `-init-token` flag.

You can get your `-calendar-id` in the preferences of your google calendar.

```
./calendar -init-token -calendar-id=mycalendar@group.calendar.google.com -event-color-id=11 -event-title="I'm home" -mac=c0:ee:bb:58:33:9e
```

```
Usage of ./calendar:
  -arp-file string
        ARP linux file (default "/proc/net/arp")
  -calendar-id string
        Google calendar id e.g. abc@group.calendar.google.com
  -client-secret-file string
        OAuth2 token file (default "client-secret.json")
  -client-token-file string
        Google OAuth2 token file, generate: -init-token (default "client-token.json")
  -event-color-id string
        Event colorid (default "11")
  -event-title string
        Event title (default "Present")
  -init-token
        Create token file
  -interval duration
        Update interval (default 1m0s)
  -mac string
        MAC to identify presence (default "00:01:02:aa:bb:cc")
```
