package resolver

type anyResolver struct {
	resolvers []Resolver
}

func NewAnyResolver(resolvers []Resolver) Resolver {
	return &anyResolver{resolvers: resolvers}
}

func (a *anyResolver) Resolve(addonName string, version string) (found bool, path string, err error) {
	for _, r := range a.resolvers {
		found, path, err = r.Resolve(addonName, version)
		if found || err != nil {
			break
		}
	}
	return
}
