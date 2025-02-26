package client

import (
	"errors"
	"slices"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
	pg "github.com/openhdc/openhdc/connectors/oracle/client/or"
	"github.com/openhdc/openhdc/metadata"
)

func toSchemaMetadata(tableName string) *arrow.Metadata {
	m := map[string]string{}
	metadata.SetTableName(m, tableName)
	md := arrow.MetadataFrom(m)
	return &md
}

func toFieldMetadata(att *pg.Attribute) arrow.Metadata {
	m := map[string]string{}
	if att.ConTypes != nil {
		if strings.Contains(*att.ConTypes, "p") {
			metadata.SetPrimaryKey(m)
		}
		if strings.Contains(*att.ConTypes, "u") {
			metadata.SetUnique(m)
		}
	}
	return arrow.MetadataFrom(m)
}

func skip(sch *arrow.Schema, keys, skips []string) bool {
	key, err := metadata.GetTableName(sch)
	if err != nil {
		return true
	}
	if slices.Contains(skips, key) {
		return true
	}
	if len(keys) > 0 && !slices.Contains(keys, key) {
		return true
	}
	return false
}

func toMessageKind(sch *arrow.Schema, mode property.SyncMode, curs []*workload.Sync_Option_Cursor) (property.MessageKind, error) {
	hasPrimaryKey := metadata.HasPrimaryKey(sch)
	hasCursor := len(curs) > 0
	switch mode {
	case property.SyncMode_full_overwrite:
		// upsert & delete
		if hasPrimaryKey {
			return property.MessageKind_upsert_update, nil
		}
		// truncate & insert
		return property.MessageKind_insert, nil

	case property.SyncMode_full_append:
		// do nothing
		if hasPrimaryKey {
			return property.MessageKind_upsert_nothing, nil
		}
		// insert
		return property.MessageKind_insert, nil

	case property.SyncMode_incremental_append:
		if !hasCursor {
			return 0, errors.New("cursor is empty")
		}
		if hasPrimaryKey {
			return 0, errors.New("primary key is exists, use `incremental_append_dedupe` instead")
		}
		// insert
		return property.MessageKind_insert, nil

	case property.SyncMode_incremental_append_dedupe:
		if !hasCursor {
			return 0, errors.New("cursor is empty")
		}
		if !hasPrimaryKey {
			return 0, errors.New("primary key is empty, use `incremental_append` instead")
		}
		// upsert_update
		return property.MessageKind_upsert_update, nil

	default:
		return property.MessageKind_message_kind_unspecified, nil
	}
}

func deleteAll(sch *arrow.Schema, syncMode property.SyncMode) bool {
	return syncMode == property.SyncMode_full_overwrite && !metadata.HasPrimaryKey(sch)
}

func deleteStale(sch *arrow.Schema, syncMode property.SyncMode) bool {
	return syncMode == property.SyncMode_full_overwrite && metadata.HasPrimaryKey(sch)
}
