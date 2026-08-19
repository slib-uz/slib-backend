package enum

type ServiceType string

const (
	ServiceTimeBased  ServiceType = "TIME_BASED"
	ServiceUsageBased ServiceType = "USAGE_BASED"
)
