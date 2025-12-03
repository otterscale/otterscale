package standalone

type base interface {
	Charms() []charm
	Config(charmName string) (string, error)
	Relations() [][]string
	Tags() []string
}
