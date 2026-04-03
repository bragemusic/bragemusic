package internalusers

import (
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

var MetaSyncer = uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))

func GetIntenalUsers() []types.UserDetails {
	return []types.UserDetails{
		{
			ID:       MetaSyncer,
			Email:    "",
			Username: "Metadata Syncer",
		},
	}
}
