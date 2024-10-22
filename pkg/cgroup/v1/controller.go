package cgroup

type Controller interface {
	Name() string
	Apply(path string) error
}
