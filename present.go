package whoshome

type PresenceProvider interface {
	Present() ([]string, error)
}
