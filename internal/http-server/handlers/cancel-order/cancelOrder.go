package cancelOrder

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
	OrderID string
}

type orderCanceller interface {
	CancelOrder(orderID string) error
}

func New(log *slog.Logger, ordCanceller orderCanceller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cancelOrder.New"

		log = log.With(slog.String("op", op))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(consts.ReqDecodeFail, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(consts.ReqDecodeFail))
			return
		}

		if err := ordCanceller.CancelOrder(req.OrderID); err != nil {
			errMsg := ""
			switch err.(type) {
			case st.ErrOrderNotExists:
				errMsg = consts.StorageInvalidOrderID
			case st.ErrHotelNotFound:
				errMsg = consts.StorageInvalidHotelID
			case st.ErrRoomNotFound:
				errMsg = consts.StorageInvalidRoomID
			case st.ErrUnavailableDate:
				errMsg = consts.StorageUnavailableDate
			default:
				errMsg = consts.StorageInternalError
			}

			log.Error(errMsg, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(errMsg))
			return
		}

		log.Info(consts.StorageOrderCancelled, req.OrderID)

		responseOK(w)
	}
}

func responseOK(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(resp.OK())
}
