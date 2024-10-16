package memory

import (
	"fmt"
	"sync"
	"time"

	st "applicationDesignTest/internal/storage"
)

type Storage struct {
	orders         map[string]st.Order
	hotelsAndRooms map[string]map[string]map[time.Time]*int
	mu             sync.RWMutex
}

func New() *Storage {
	return &Storage{
		orders:         make(map[string]st.Order),
		hotelsAndRooms: make(map[string]map[string]map[time.Time]*int),
		mu:             sync.RWMutex{},
		//roomAvailability: make([]RoomAvailability),
	}
}

func (s *Storage) UpdateInventory(roomAvailability ...st.RoomAvailability) error {
	const op = "storage.memory.UpdateRooms"

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range roomAvailability {
		r := r
		//s.roomAvailability = append(s.roomAvailability, roomRecord)

		if _, ok := s.hotelsAndRooms[r.HotelID]; !ok {
			s.hotelsAndRooms[r.HotelID] = make(map[string]map[time.Time]*int)
		}

		if _, ok := s.hotelsAndRooms[r.HotelID][r.RoomID]; !ok {
			s.hotelsAndRooms[r.HotelID][r.RoomID] = make(map[time.Time]*int)
		}

		s.hotelsAndRooms[r.HotelID][r.RoomID][toDay(r.Date)] = &r.Quota
	}

	return nil
}

func (s *Storage) GetHotels() ([]st.RoomAvailability, error) {
	const op = "storage.memory.GetHotels"

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []st.RoomAvailability
	for hotelID, roomsMap := range s.hotelsAndRooms {
		for roomID, roomSchedule := range roomsMap {
			for date, quota := range roomSchedule {
				r := st.RoomAvailability{
					HotelID: hotelID,
					RoomID:  roomID,
					Date:    date,
					Quota:   *quota,
				}
				result = append(result, r)
			}
		}
	}

	return result, nil
}

func (s *Storage) GetRooms(hotelID string) ([]st.RoomAvailability, error) {
	const op = "storage.memory.GetRooms"

	s.mu.RLock()
	defer s.mu.RUnlock()

	roomsMap, ok := s.hotelsAndRooms[hotelID]
	if !ok {
		err := st.ErrHotelNotFound{HotelID: hotelID}
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	var result []st.RoomAvailability
	for roomID, roomSchedule := range roomsMap {
		for date, quota := range roomSchedule {
			r := st.RoomAvailability{
				HotelID: hotelID,
				RoomID:  roomID,
				Date:    date,
				Quota:   *quota,
			}
			result = append(result, r)
		}
	}

	return result, nil
}

//func (s *Storage) GetOrders() ([]st.Order, error) {
//	const op = "storage.memory.GetOrders"
//
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	var result []st.Order
//	for _, order := range s.orders {
//		result = append(result, order)
//	}
//
//	return result, nil
//}

func (s *Storage) CreateOrder(order st.Order) error {
	const op = "storage.memory.CreateOrder"

	s.mu.RLock()
	_, ok := s.orders[order.OrderID]
	if ok {
		err := st.ErrOrderExists{OrderID: order.OrderID}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	roomsMap, ok := s.hotelsAndRooms[order.HotelID]
	if !ok {
		err := st.ErrHotelNotFound{HotelID: order.HotelID}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	roomSchedule, ok := roomsMap[order.RoomID]
	if !ok {
		err := st.ErrRoomNotFound{HotelID: order.HotelID, RoomID: order.RoomID}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	s.mu.RUnlock()

	nDays, err := daysBetween(order.From, order.To)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// another approach -> atomics
	s.mu.Lock()
	defer s.mu.Unlock()

	var fetchedQuotas []*int
	curDate := toDay(order.From)

	for i := 0; i <= nDays; i++ {
		dailyQuota, ok := roomSchedule[curDate]
		if !ok {
			err := st.ErrUnavailableDate{HotelID: order.HotelID, RoomID: order.RoomID, Date: curDate}
			return fmt.Errorf("%s: execute statement: %w", op, err)
		}

		if *dailyQuota < 1 {
			err := st.ErrNoRoomsForDate{HotelID: order.HotelID, RoomID: order.RoomID, Date: curDate}
			return fmt.Errorf("%s: execute statement: %w", op, err)
		}

		fetchedQuotas = append(fetchedQuotas, dailyQuota)
		curDate = curDate.AddDate(0, 0, 1)
	}

	for _, dailyQuota := range fetchedQuotas {
		*dailyQuota--
	}

	order.Status = st.StatusReserved
	s.orders[order.OrderID] = order

	return nil
}

func (s *Storage) CancelOrder(orderID string) error {
	const op = "storage.memory.CancelOrder"

	s.mu.RLock()
	order, ok := s.orders[orderID]
	if !ok {
		err := st.ErrOrderNotExists{OrderID: orderID}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	roomsMap, ok := s.hotelsAndRooms[order.HotelID]
	if !ok {
		err := st.ErrHotelNotFound{HotelID: order.HotelID}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	roomSchedule, ok := roomsMap[order.RoomID]
	if !ok {
		err := st.ErrRoomNotFound{HotelID: order.HotelID, RoomID: order.RoomID}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	s.mu.RUnlock()

	nDays, err := daysBetween(order.From, order.To)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var fetchedQuotas []*int
	curDate := toDay(order.From)

	for i := 0; i <= nDays; i++ {
		dailyQuota, ok := roomSchedule[curDate]
		if !ok {
			err := st.ErrUnavailableDate{HotelID: order.HotelID, RoomID: order.RoomID, Date: curDate}
			return fmt.Errorf("%s: execute statement: %w", op, err)
		}

		fetchedQuotas = append(fetchedQuotas, dailyQuota)
		curDate = curDate.AddDate(0, 0, 1)
	}

	for _, dailyQuota := range fetchedQuotas {
		*dailyQuota++
	}

	order.Status = st.StatusCancelled
	s.orders[order.OrderID] = order

	return nil
}
