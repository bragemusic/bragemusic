package audiointerface

import (
	"context"

	"github.com/bragemusic/core/pkg/types"
)

type AudioInterface interface {
	IsPlaying() bool
	PlayedMS() int64
	Pause()
	Play()
	PlayPause()
	RegisterErrorCallback(f func(context.Context, error))
	StartAudioFile(ctx context.Context, audioFile types.AudioFile, finishedCallback func(context.Context) error)
	Stop()
	Terminate()
}
