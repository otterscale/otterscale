package workload

import "github.com/openhdc/openhdc/api/property/v1"

func GetSyncMode(os []*Sync_Option, v string) property.SyncMode {
	if v == "" {
		return 0
	}
	for _, o := range os {
		if o.GetKey() == v {
			return o.GetMode()
		}
	}
	return 0
}

func GetSyncCursor(os []*Sync_Option, v string) string {
	if v == "" {
		return ""
	}
	for _, o := range os {
		if o.GetKey() == v {
			return o.GetCursor()
		}
	}
	return ""
}
