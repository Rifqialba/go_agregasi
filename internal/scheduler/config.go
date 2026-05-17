package scheduler

type PollingSourceConfig struct {
	SourceID string `yaml:"source_id"`
	URL      string `yaml:"url"`
	Schedule string `yaml:"schedule"`
}

type SchedulerConfig struct {
	PollingSources []PollingSourceConfig `yaml:"polling_sources"`
}