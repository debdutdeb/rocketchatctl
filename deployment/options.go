package deployment

type HostDeploymentOptions struct {
	RootUrl           string
	Version           string
	MongoVersion      string
	RegistrationToken string
	Port              int16
	BindLoopback      bool
	UseExistingMongo  bool
}
