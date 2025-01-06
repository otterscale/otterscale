package pb

import "github.com/openhdc/openhdc/api/property/v1"

const (
	Migrate       = property.MessageKind_migrate
	Insert        = property.MessageKind_insert
	UpsertUpdate  = property.MessageKind_upsert_update
	UpsertNothing = property.MessageKind_upsert_nothing
	DeleteStale   = property.MessageKind_delete_stale
	DeleteAll     = property.MessageKind_delete_all
)
