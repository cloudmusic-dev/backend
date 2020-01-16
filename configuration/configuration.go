package configuration

type Database struct {
	Host     string
	Username string
	Password string
	Database string
}

type Configuration struct {
	Database Database
}
