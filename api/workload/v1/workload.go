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

func GetSyncCursors(os []*Sync_Option, v string) []*Sync_Option_Cursor {
	if v == "" {
		return nil
	}
	for _, o := range os {
		if o.GetKey() == v {
			return o.GetCursors()
		}
	}
	return nil
}
