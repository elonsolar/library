package config

// Config for trace
type Config struct {
	ResourceConfig *ResourceConfig
	ExportConfig   *ExportConfig
}

//  config for resouce ,todo
type ResourceConfig struct {
	ServiceName    string
	ServiceVersion string
}
