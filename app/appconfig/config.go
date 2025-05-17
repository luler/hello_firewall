package appconfig

type Config struct {
	App struct {
		Env string
	}
	Database map[string]struct {
		Driver       string
		Name         string
		Host         string
		Port         int
		Username     string
		Password     string
		Table_Prefix string
	}
	Redis map[string]struct {
		Host     string
		Port     int
		Password string
		Select   int
	}
}
