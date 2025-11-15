package standalone

type base interface {
	Charms() []charm
	Configs() (string, error)
	Relations() [][]string
	Tags() []string
}
