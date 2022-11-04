package reports

type Config struct {
	RootPath   string
	APIVersion string
}

func (c *Config) pathPrefix() string {
	return c.RootPath + c.APIVersion
}

var config = Config{
	RootPath:   "/reports/",
	APIVersion: "2021-06-30",
}
