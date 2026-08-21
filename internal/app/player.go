package app

import (
	"errors"
	"fmt"

	"github.com/bragemusic/core/pkg/database"
	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
)

func (a *App) Start(contextType types.PlayContextType, parentID string, trackNumber int) {
	uid, err := uuid.FromString(parentID)
	if err != nil {
		a.handleError(err)
		return
	}

	switch contextType {
	case types.PlayContextAlbum:
		err = a.client.StartPlayerWithAlbum(a.ctx, uid, trackNumber)
		if err != nil {
			a.handleError(err)
			return
		}
	case types.PlayContextPlaylist:
		err = a.client.StartPlayerWithPlaylist(a.ctx, uid, trackNumber, database.SortByDate, database.SortAsc)
		if err != nil {
			a.handleError(err)
			return
		}
	case types.PlayContextLikedTracks:
		err = a.client.StartPlayerWithLikedTracks(a.ctx, trackNumber)
		if err != nil {
			a.handleError(err)
			return
		}
	// case types.PlayContextFilter:
	// err = a.client.StartPlayerWithTrackFilter(a.ctx, *filter, trackNumber)
	// 	if err != nil {
	// 		a.handleError(err)
	// 		return
	// 	}
	default:
		a.handleError(errors.New("unknown audioplayer context type"))
		return
	}

	a.SendMessage("started")
}

func (a *App) NextRepeat() {
	cr := a.client.PlayerState().Playback.Repeat
	switch cr {
	case types.RepeatOff:
		a.client.SetRepeat(a.ctx, types.RepeatAll)

	case types.RepeatAll:
		a.client.SetRepeat(a.ctx, types.RepeatOne)

	case types.RepeatOne:
		a.client.SetRepeat(a.ctx, types.RepeatOff)
	}
}

func (a *App) ToggleShuffle() {
	cs := a.client.PlayerState().Playback.Shuffle
	a.client.SetShuffle(a.ctx, !cs)
}

func (a *App) PlayPause() {
	a.client.PlayPause(a.ctx)
	a.SendMessage("paus/play")
}

func (a *App) NextTrack() {
	err := a.client.NextTrack(a.ctx)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) PreviousTrack() {
	err := a.client.PreviousTrack(a.ctx)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) AddTrackToQueue(trackID, albumID string) {
	tID, err := uuid.FromString(trackID)
	if err != nil {
		a.handleError(err)
		return
	}

	aID, err := uuid.FromString(albumID)
	if err != nil {
		a.handleError(err)
		return
	}

	fmt.Println(trackID, albumID)

	err = a.client.AddTrackToQueue(a.ctx, tID, aID)
	if err != nil {
		a.handleError(err)
		return
	}
}
