package property

func ParseWorkloadKind(str string) WorkloadKind {
	if v, ok := WorkloadKind_value[str]; ok {
		return WorkloadKind(v)
	}
	return WorkloadKind(0)
}

func ParseSyncMode(str string) SyncMode {
	if v, ok := SyncMode_value[str]; ok {
		return SyncMode(v)
	}
	return SyncMode(0)
}
