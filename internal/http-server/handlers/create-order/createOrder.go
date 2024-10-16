package createOrder

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
	Order st.Order
}

type orderCreator interface {
	CreateOrder(order st.Order) error
}

func New(log *slog.Logger, ordCreator orderCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.createOrder.New"

		log = log.With(slog.String("op", op))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(consts.ReqDecodeFail, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(consts.ReqDecodeFail))
			return
		}

		if err := ordCreator.CreateOrder(req.Order); err != nil {
			errMsg := ""
			switch err.(type) {
			case st.ErrOrderExists:
				errMsg = consts.StorageExistedOrderID
			case st.ErrHotelNotFound:
				errMsg = consts.StorageInvalidHotelID
			case st.ErrRoomNotFound:
				errMsg = consts.StorageInvalidRoomID
			case st.ErrUnavailableDate:
				errMsg = consts.StorageUnavailableDate
			case st.ErrNoRoomsForDate:
				errMsg = consts.StorageNoRoomsForPeriod
			default:
				errMsg = consts.StorageInternalError
			}

			log.Error(errMsg, sl.Err(err))
			json.NewEncoder(w).Encode(resp.Error(errMsg))
			return
		}

		log.Info(consts.StorageOrderCreated, req.Order.OrderID)

		responseOK(w)
	}
}

func responseOK(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(resp.OK())
}
