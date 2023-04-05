package common

// ResourcesContainerInf is the interface common to all resources containers
type ResourcesContainerInf interface {
	CollectResourcesFromAPI() error
	PrintStats()
	ToString() (string, error)
}
