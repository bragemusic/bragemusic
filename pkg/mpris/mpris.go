package mpris

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bragemusic/core/pkg/types"
	"github.com/godbus/dbus/v5"
)

type MprisStatus string

const (
	MprisPlaying MprisStatus = "Playing"
	MprisPaused  MprisStatus = "Paused"
)

// --- org.mpris.MediaPlayer2 interface ---
type MediaPlayer2 struct {
	playerName string
}

func (m *MediaPlayer2) Raise() *dbus.Error { return nil }
func (m *MediaPlayer2) Quit() *dbus.Error  { return nil }
func (m *MediaPlayer2) Get() *dbus.Error   { return nil }
func (m *MediaPlayer2) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	return map[string]dbus.Variant{
		"CanQuit":      dbus.MakeVariant(true),
		"CanRaise":     dbus.MakeVariant(false),
		"HasTrackList": dbus.MakeVariant(false),
		"Identity":     dbus.MakeVariant(m.playerName),
		"DesktopEntry": dbus.MakeVariant(m.playerName),
	}, nil
}

// --- org.mpris.MediaPlayer2.Player interface ---
type Player struct {
	mu              sync.Mutex
	status          MprisStatus
	conn            *dbus.Conn
	playerPlay      func(context.Context)
	playerPause     func(context.Context)
	playerPlayPause func(context.Context)
	playerNext      func(context.Context) error
	playerPrevious  func(context.Context) error
	track           *types.TrackDetailed
	ctx             context.Context
}

func (p *Player) Metadata() map[string]dbus.Variant {
	if p.track == nil {
		return nil
	}

	return map[string]dbus.Variant{
		"xesam:album":       dbus.MakeVariant(p.track.AlbumName),
		"xesam:artist":      dbus.MakeVariant(p.track.ArtistNames),
		"xesam:discNumber":  dbus.MakeVariant(p.track.DiscNumber),
		"xesam:title":       dbus.MakeVariant(p.track.Title),
		"xesam:trackNumber": dbus.MakeVariant(p.track.TrackNumber),
	}
}

func (p *Player) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return map[string]dbus.Variant{
		"PlaybackStatus": dbus.MakeVariant(p.status),
		"Metadata":       dbus.MakeVariant(p.Metadata()),
		"CanGoNext":      dbus.MakeVariant(true),
		"CanGoPrevious":  dbus.MakeVariant(true),
		"CanPlay":        dbus.MakeVariant(true),
		"CanPause":       dbus.MakeVariant(true),
		"CanSeek":        dbus.MakeVariant(false),
		"CanControl":     dbus.MakeVariant(true),
		"Shuffle":        dbus.MakeVariant(false),
		"LoopStatus":     dbus.MakeVariant("None"),
		"Rate":           dbus.MakeVariant(1.0),
		"Volume":         dbus.MakeVariant(1.0),
		"MinimumRate":    dbus.MakeVariant(1.0),
		"MaximumRate":    dbus.MakeVariant(1.0),
	}, nil
}

func (p *Player) Get(iface, prop string) (dbus.Variant, *dbus.Error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch prop {
	case "PlaybackStatus":
		return dbus.MakeVariant(p.status), nil
	case "Metadata":
		return dbus.MakeVariant(p.Metadata()), nil
	}
	return dbus.Variant{}, nil
}

func (p *Player) Set(iface, prop string, value dbus.Variant) *dbus.Error {
	return nil
}

// --- playback methods ---
func (p *Player) Play() *dbus.Error {
	p.playerPlay(p.ctx)
	return nil
}

func (p *Player) Pause() *dbus.Error {
	p.playerPause(p.ctx)
	return nil
}

func (p *Player) Next() *dbus.Error {
	p.playerNext(p.ctx)
	return nil
}

func (p *Player) Previous() *dbus.Error {
	p.playerPrevious(p.ctx)
	return nil
}

func (p *Player) PlayPause() *dbus.Error {
	p.playerPlayPause(p.ctx)
	return nil
}

func (p *Player) Stop() *dbus.Error {
	return nil
}

func (p *Player) SetStatus(s MprisStatus) {
	p.mu.Lock()
	p.status = s
	p.conn.Emit("/org/mpris/MediaPlayer2",
		"org.freedesktop.DBus.Properties.PropertiesChanged",
		"org.mpris.MediaPlayer2.Player",
		map[string]dbus.Variant{"PlaybackStatus": dbus.MakeVariant(s)},
		[]string{},
	)
	p.mu.Unlock()
}

func (p *Player) SetTrack(track *types.TrackDetailed) {
	p.mu.Lock()
	p.track = track
	p.conn.Emit("/org/mpris/MediaPlayer2",
		"org.freedesktop.DBus.Properties.PropertiesChanged",
		"org.mpris.MediaPlayer2.Player",
		map[string]dbus.Variant{"Metadata": dbus.MakeVariant(p.Metadata())},
		[]string{},
	)
	p.mu.Unlock()
}

type Mpris struct {
	p          *Player
	mp2        *MediaPlayer2
	conn       *dbus.Conn
	playerName string
}

func (m *Mpris) SetStatus(s MprisStatus) {
	m.p.SetStatus(s)
}

func (m *Mpris) SetTrack(track *types.TrackDetailed) {
	m.p.SetTrack(track)
}

func New(playerName string, playFunc, pauseFunc, playPauseFunc func(ctx context.Context), prevFunc, nextFunc func(context.Context) error) (mp Mpris, err error) {
	mp = Mpris{
		p: &Player{
			playerPlay:      playFunc,
			playerPlayPause: playPauseFunc,
			status:          "Paused",
			playerPause:     pauseFunc,
			playerNext:      nextFunc,
			playerPrevious:  prevFunc,
		},
		mp2: &MediaPlayer2{
			playerName: playerName,
		},
		playerName: playerName,
	}

	mp.conn, err = dbus.SessionBus()
	if err != nil {
		return Mpris{}, err
	}

	reply, err := mp.conn.RequestName(fmt.Sprintf("org.mpris.MediaPlayer2.%s", mp.playerName), dbus.NameFlagDoNotQueue)
	if err != nil {
		return Mpris{}, err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return Mpris{}, errors.New("name already taken")
	}

	mp.p.conn = mp.conn

	// Export separate interfaces explicitly
	if err := mp.conn.Export(mp.mp2, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2"); err != nil {
		return Mpris{}, err
	}

	if err := mp.conn.Export(mp.p, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player"); err != nil {
		return Mpris{}, err
	}

	if err := mp.conn.Export(mp.p, "/org/mpris/MediaPlayer2", "org.freedesktop.DBus.Properties"); err != nil {
		return Mpris{}, err
	}

	return mp, nil
}
