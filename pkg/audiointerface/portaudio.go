package audiointerface

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"log/slog"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gordonklaus/portaudio"
)

type PortAudio struct {
	stream       *portaudio.Stream
	audioFile    types.AudioFile
	isPlaying    bool
	framesPlayed int64
	log          *slog.Logger

	errCallback func(context.Context, error)

	cStop  chan bool
	cPlay  chan bool
	cPause chan bool
}

func (p *PortAudio) StartAudioFile(ctx context.Context, audioFile types.AudioFile, finishedCallback func(context.Context) error) {
	p.audioFile = audioFile
	p.killStream()

	p.isPlaying = true
	go p.runStream(ctx, finishedCallback)
}

func (p *PortAudio) RegisterErrorCallback(f func(context.Context, error)) {
	p.errCallback = f
}

func (p *PortAudio) runStream(ctx context.Context, finishedCallback func(context.Context) error) {
	out := make([]float32, 8192) // buffer of float32 samples

	var err error
	p.stream, err = portaudio.OpenDefaultStream(0, p.audioFile.NChannels(), float64(p.audioFile.SampleRate()), len(out)/p.audioFile.NChannels(), &out)
	if err != nil {
		p.handleError(ctx, err)
		return
	}

	p.log.DebugContext(ctx, "starting new stream")

	err = p.stream.Start()
	if err != nil {
		p.handleError(ctx, err)
		return
	}

	p.framesPlayed = 0

	for {
		p.isPlaying = true

		audio := make([]byte, 4*len(out)) // 4 bytes per float32
		n, err := p.audioFile.Read(audio)
		if err != nil {
			if err == io.EOF && n == 0 {
				break
			} else {
				p.handleError(ctx, err)
				return
			}
		}

		frames := n / (4 * p.audioFile.NChannels())
		p.framesPlayed += int64(frames)

		// Convert bytes to float32 slice
		err = binary.Read(bytes.NewReader(audio[:n]), binary.LittleEndian, out[:n/4])
		if err != nil {
			p.handleError(ctx, err)
			return
		}

		err = p.stream.Write()
		if err != nil {
			p.handleError(ctx, err)
			return
		}

		// Stop/pause handling
		select {
		case <-p.cStop:
			p.killStream()
			p.log.DebugContext(ctx, "stopping stream")
			return
		case <-p.cPause:
			p.log.DebugContext(ctx, "pausing stream")
			err = p.stream.Stop()
			if err != nil {
				p.handleError(ctx, err)
				return
			}

			select {
			case <-p.cPlay:
				p.log.DebugContext(ctx, "playing stream")
				err = p.stream.Start()
				if err != nil {
					p.handleError(ctx, err)
					return
				}

			case <-p.cStop:
				p.killStream()
				p.log.DebugContext(ctx, "stopping stream")
				return
			}
			p.log.DebugContext(ctx, "playing stream")
		default:
		}
	}

	p.killStream()
	finishedCallback(ctx)
}

func (p *PortAudio) IsPlaying() bool {
	return p.isPlaying
}

func (p *PortAudio) PlayedMS() int64 {
	if p.audioFile == nil {
		return 0
	}
	return int64(float64(p.framesPlayed) / float64(p.audioFile.SampleRate()) * 1000.0)
}

func (p *PortAudio) Stop() {
	if p.stream != nil {
		p.cStop <- true
	}
	p.isPlaying = false
	// p.stream = nil
}

func (p *PortAudio) PlayPause() {
	if p.stream == nil {
		return
	}

	if p.isPlaying {
		p.Pause()
	} else {
		p.Play()
	}
}

func (p *PortAudio) Pause() {
	p.cPause <- true
	p.isPlaying = false
}

func (p *PortAudio) Play() {
	p.cPlay <- true
	p.isPlaying = true
}

func (p *PortAudio) killStream() {
	if p.stream != nil {
		p.stream.Stop()
		p.stream.Close()
		p.stream = nil
	}
	p.isPlaying = false
}

func (p *PortAudio) Terminate() {
	p.killStream()
	portaudio.Terminate()
}

func (p *PortAudio) handleError(ctx context.Context, err error) {
	p.log.ErrorContext(ctx, err.Error())
	if p.errCallback != nil {
		p.errCallback(ctx, err)
	}
	p.killStream()
}

func NewPortAudio(slogHandler slog.Handler) (*PortAudio, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}

	return &PortAudio{
		cStop:  make(chan bool),
		cPlay:  make(chan bool),
		cPause: make(chan bool),
		log:    slog.New(slogHandler).With("service", "portaudio"),
	}, nil
}
