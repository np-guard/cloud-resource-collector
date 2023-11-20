package common

// ResourcesContainerInf is the interface common to all resources containers
type ResourcesContainerInf interface {
	CollectResourcesFromAPI() error
	PrintStats()
	ToJSONString() (string, error)
	AllRegions() []string
}
