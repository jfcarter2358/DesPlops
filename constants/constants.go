package constants

const VERSION = "0.1.0"

// k8s manifest kinds that cause kapp to break
// so we want to filter them out
var API_FILTER = []string{
	"Endpoints",
	"ReplicaSet",
	"Pod",
	"EndpointSlice",
	"Namespace",
	"CustomResourceDefinition",
}
