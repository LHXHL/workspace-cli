package models

// Config holds the product-specific configuration.
type Config struct {
	URL      string `yaml:"url"`
	APIToken string `yaml:"api_token"`
}
