package workload

type Kind int

const (
	KindSource Kind = iota
	KindDestination
	KindTransformer
)

func (k Kind) String() string {
	switch k {
	case KindSource:
		return "source"
	case KindDestination:
		return "destination"
	case KindTransformer:
		return "transformer"
	}
	return ""
}
