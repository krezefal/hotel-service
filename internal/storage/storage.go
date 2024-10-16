package storage

import (
	"fmt"
	"time"
)

const (
	StatusReserved  = "reserved"
	StatusCancelled = "cancelled"
)

type Order struct {
	OrderID   string    `json:"order_id"`
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	Status    string    `json:"status"`
}

type RoomAvailability struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	Date    time.Time `json:"date"`
	Quota   int       `json:"quota"`
}

type ErrOrderExists struct {
	OrderID string
}

func (e ErrOrderExists) Error() string {
	return fmt.Sprintf("the order with id=%s already exists", e.OrderID)
}

type ErrOrderNotExists struct {
	OrderID string
}

func (e ErrOrderNotExists) Error() string {
	return fmt.Sprintf("order with id=%s does not exist", e.OrderID)
}

type ErrHotelNotFound struct {
	HotelID string
}

func (e ErrHotelNotFound) Error() string {
	return fmt.Sprintf("hotel not found by provided id: %s", e.HotelID)
}

type ErrRoomNotFound struct {
	HotelID string
	RoomID  string
}

func (e ErrRoomNotFound) Error() string {
	return fmt.Sprintf("hotel id=%s: room not found by provided id: %s", e.HotelID, e.RoomID)
}

type ErrUnavailableDate struct {
	HotelID string
	RoomID  string
	Date    time.Time
}

func (e ErrUnavailableDate) Error() string {
	return fmt.Sprintf("hotel id=%s, room id=%s: unavailable date: %s", e.HotelID, e.RoomID, e.Date)
}

type ErrNoRoomsForDate struct {
	HotelID string
	RoomID  string
	Date    time.Time
}

func (e ErrNoRoomsForDate) Error() string {
	return fmt.Sprintf("hotel id=%s, room id=%s: no rooms available for date: %s", e.HotelID, e.RoomID, e.Date)
}
