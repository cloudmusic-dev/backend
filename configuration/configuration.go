package configuration

type Database struct {
	Host     string
	Username string
	Password string
	Database string
}

type Signing struct {
	Method         string
	Strength       int
	Key            string
	PrivateKeyPath string
	PublicKeyPath  string
}

type Configuration struct {
	Database Database
	Signing Signing
}
