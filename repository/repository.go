package repository

type Repository interface {
	Store(addonPath string, addonName string, version string) error
}
