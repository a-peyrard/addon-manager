package resolver

type Resolver interface {
	Resolve(packageName string, version string) (found bool, path string, err error)
}
