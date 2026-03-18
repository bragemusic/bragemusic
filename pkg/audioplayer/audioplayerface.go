package audioplayer

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

type AudioPlayerFace interface {
	RegisterErrorCallback(f func(context.Context, error))
	RegisterPlayContextCallback(f func(context.Context, types.PlayContext))
	RegisterPlaybackStateCallback(f func(context.Context, types.PlaybackState))
	RegisterPlayCountCallback(f func(trackID uuid.UUID))
	PlayerState() types.PlayerState
	LoadAndStartTracks(ctx context.Context, state types.PlayerState) (err error)
	SetRepeat(ctx context.Context, r types.RepeatType)
	SetShuffle(ctx context.Context, s bool)
	NextTrack(ctx context.Context) (err error)
	PreviousTrack(ctx context.Context) (err error)
	Pause(ctx context.Context)
	Play(ctx context.Context)
	PlayPause(ctx context.Context)
	Stop(ctx context.Context) error
	AddTrackToQueue(ctx context.Context, track types.TrackDetailed)
}
