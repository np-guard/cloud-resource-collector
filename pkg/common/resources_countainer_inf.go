package common

type ResourcesContainerInf interface {
	CollectResourcesFromAPI()
	PrintStats()
	ToString() (string, error)
}
