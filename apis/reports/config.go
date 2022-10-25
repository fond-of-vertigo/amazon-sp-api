package reports

type Config struct {
	RootRoute  string
	APIVersion string
}

func (c *Config) routePrefix() string {
	return c.RootRoute + c.APIVersion
}

var config = Config{
	RootRoute:  "/reports/",
	APIVersion: "2021-06-30",
}
