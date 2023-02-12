package resolver

type githubResolver struct{}

func NewGithubResolver() Resolver {
	return &githubResolver{}
}

func (g *githubResolver) Resolve(packageName string, version string) (found bool, path string, err error) {
	//TODO implement me
	panic("implement me")
}
