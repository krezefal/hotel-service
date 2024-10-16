package updateInventory

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"applicationDesignTest/internal/consts"
	resp "applicationDesignTest/internal/lib/api/response"
	sl "applicationDesignTest/internal/lib/logger/slog"
	st "applicationDesignTest/internal/storage"
)

type Request struct {
	RoomAvailability []st.RoomAvailability
}

type InventoryUpdater interface {
	UpdateInventory(roomAvailability ...st.RoomAvailability) error
}

func New(log *slog.Logger, invUpdater InventoryUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.updateInventory.New"

		log = log.With(slog.String("op", op))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(consts.ReqDecodeFail, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(consts.ReqDecodeFail))
			return
		}

		if err := invUpdater.UpdateInventory(req.RoomAvailability...); err != nil {
			log.Error(consts.StorageInternalError, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(consts.StorageInternalError))
			return
		}

		log.Info(consts.StorageInventoryUpdated)

		responseOK(w)
	}
}

func responseOK(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(resp.OK())
}
