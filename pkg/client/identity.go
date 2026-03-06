package client

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type Identity struct {
	id   uuid.UUID
	name string
}

func (i Identity) ClientID() uuid.UUID {
	return i.id
}

func (i Identity) ClientName() string {
	return i.name
}

func (i Identity) ClientType() types.DeviceType {
	return types.DeviceTypeStreaming
}

func idFilePath() (string, error) {
	return xdg.StateFile(filepath.Join("brage", "client_id"))
}

func NewIdentity(name string) (Identity, error) {
	i := Identity{
		name: name,
	}
	path, err := idFilePath()
	if err != nil {
		return Identity{}, err
	}

	b, err := os.ReadFile(path)
	// Happy path
	if err == nil {
		uid, err := uuid.FromString(string(b))
		if err == nil {
			i.id = uid
			return i, nil
		}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Identity{}, err
	}

	if err := os.WriteFile(path, []byte(uid.String()), 0o600); err != nil {
		return Identity{}, nil
	}

	i.id = uid

	return i, nil
}
