// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"applicationDesignTest/internal/config"
	sl "applicationDesignTest/internal/lib/logger/slog"
	st "applicationDesignTest/internal/storage"
	"applicationDesignTest/internal/storage/memory"

	"applicationDesignTest/internal/http-server/handlers/cancel-order"
	"applicationDesignTest/internal/http-server/handlers/create-order"
	"applicationDesignTest/internal/http-server/handlers/get-hotels"
	"applicationDesignTest/internal/http-server/handlers/get-rooms"
	"applicationDesignTest/internal/http-server/handlers/update-inventory"
)

var Availability = []st.RoomAvailability{
	{"reddison", "lux", date(2024, 1, 1), 1000},
	{"reddison", "lux", date(2024, 1, 2), 1000},
	{"reddison", "lux", date(2024, 1, 3), 1000},
	{"reddison", "lux", date(2024, 1, 4), 500},
	{"reddison", "lux", date(2024, 1, 5), 0},
	{"azimut", "lux", date(2024, 1, 5), 5},
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := sl.SetupLogger(cfg.Env)

	log.Info("Started Booking service", slog.String("env", cfg.Env))
	log.Debug("Enabled debug messages")

	storage := memory.New()
	if err := storage.UpdateInventory(Availability...); err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	router := chi.NewRouter()
	//router.Use(middleware.Recoverer)

	router.Get("/hotels", getHotels.New(log, storage))
	router.Get("/hotels/{hotel_id}", getRooms.New(log, storage))
	router.Put("/hotels/update", updateInventory.New(log, storage))
	router.Post("/hotels/order", createOrder.New(log, storage))
	router.Put("/hotels/order/&id", cancelOrder.New(log, storage))

	log.Info("Starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("Server started")

	<-done
	log.Info("Server stopped")
}
