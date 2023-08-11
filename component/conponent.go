package component

type Component interface {
	Name() string
	Registry() error
}
