package datastores

type Preferences struct {
	ServerURL      string
	ServerUsername string
	ServerPassword string
	HistorySize    int
}

type Datastore interface {
	InitDatastore() error
	LoadPreferences() (Preferences, error)
	SavePreferences(Preferences) error
}
