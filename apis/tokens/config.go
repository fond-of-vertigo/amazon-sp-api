package tokens

type Config struct {
	RootPath   string
	APIVersion string
}

func (c *Config) pathPrefix() string {
	return c.RootPath + c.APIVersion
}

var config = Config{
	RootPath:   "/tokens/",
	APIVersion: "2021-03-01",
}
