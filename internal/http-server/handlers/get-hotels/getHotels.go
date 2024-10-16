package getHotels

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"applicationDesignTest/internal/consts"
	resp "applicationDesignTest/internal/lib/api/response"
	sl "applicationDesignTest/internal/lib/logger/slog"
	st "applicationDesignTest/internal/storage"
)

type Response struct {
	resp.Response
	RoomAvailability []st.RoomAvailability `json:"alias,omitempty"`
}

type HotelsProvider interface {
	GetHotels() ([]st.RoomAvailability, error)
}

func New(log *slog.Logger, hProvider HotelsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getHotels.New"

		log = log.With(slog.String("op", op))

		roomAvailability, err := hProvider.GetHotels()
		if err != nil {
			log.Error(consts.StorageInternalError, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(consts.StorageInternalError))
			return
		}

		log.Info(consts.StorageDataFetched)

		responseOK(w, roomAvailability)
	}
}

func responseOK(w http.ResponseWriter, roomAvailability []st.RoomAvailability) {
	response := Response{
		Response:         resp.OK(),
		RoomAvailability: roomAvailability,
	}

	json.NewEncoder(w).Encode(response)
}
