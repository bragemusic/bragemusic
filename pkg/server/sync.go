package server

import (
	"encoding/json"
	"net/http"
)

func (s Server) sync() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		syncReq := SyncReq{}
		if err := json.NewDecoder(r.Body).Decode(&syncReq); err != nil {
			return http.StatusBadRequest, nil, err
		}

		syncState, err := s.mediamgr.GetSyncState(ctx, syncReq.ChangesSince)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		return http.StatusOK, syncState, nil
	},
	)
}

func (s Server) syncPlayHistory() http.HandlerFunc {
	return s.handleJSON(func(w http.ResponseWriter, r *http.Request) (int, any, error) {
		ctx := r.Context()

		syncReq := SyncPlayHistoryReq{}
		if err := json.NewDecoder(r.Body).Decode(&syncReq); err != nil {
			return http.StatusBadRequest, nil, err
		}

		syncState, err := s.mediamgr.SyncPlayHistory(ctx, syncReq.ChangesSince, syncReq.UpdatedClientItems)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		return http.StatusOK, syncState, nil
	},
	)
}
