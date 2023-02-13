package resolver

type Resolver interface {
	Resolve(addonName string, version string) (found bool, path string, err error)
}
