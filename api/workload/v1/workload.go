package workload

func ParseKind(str string) Kind {
	if v, ok := Kind_value[str]; ok {
		return Kind(v)
	}
	return Kind(0)
}
