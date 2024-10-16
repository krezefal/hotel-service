package getRooms

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"

	"applicationDesignTest/internal/consts"
	resp "applicationDesignTest/internal/lib/api/response"
	sl "applicationDesignTest/internal/lib/logger/slog"
	st "applicationDesignTest/internal/storage"
)

//type Request struct {
//	HotelID string
//}

type Response struct {
	resp.Response
	RoomAvailability []st.RoomAvailability `json:"alias,omitempty"`
}

type RoomsProvider interface {
	GetRooms(hotelID string) ([]st.RoomAvailability, error)
}

func New(log *slog.Logger, rProvider RoomsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getRooms.New"

		log = log.With(slog.String("op", op))

		hotelID := chi.URLParam(r, "hotel_id")

		//var req Request
		//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//	log.Error(consts.ReqDecodeFail, sl.Err(err))
		//	json.NewEncoder(w).Encode(resp.Error(consts.ReqDecodeFail))
		//	return
		//}

		roomAvailability, err := rProvider.GetRooms(hotelID)
		if err != nil {
			errMsg := ""
			switch err.(type) {
			case st.ErrHotelNotFound:
				errMsg = consts.StorageInvalidHotelID
			default:
				errMsg = consts.StorageInternalError
			}

			log.Error(errMsg, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(errMsg))
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
