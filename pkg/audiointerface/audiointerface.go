package audiointerface

import (
	"context"

	"github.com/bragemusic/bragemusic/pkg/types"
)

type AudioInterface interface {
	IsPlaying() bool
	PlayedMS() int64
	Pause()
	Play()
	PlayPause()
	RegisterErrorCallback(f func(context.Context, error))
	StartAudioFile(ctx context.Context, audioFile types.AudioFile, finishedCallback, timeoutCallback func(context.Context) error)
	Stop()
	Terminate()
}
