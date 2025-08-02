package order

import (
	"sync"
	"time"
)

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	ID        string    `json:"id"`
	Customer  string    `json:"customer"`
	Items     []Item    `json:"items"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type OrderStore struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[string]Order),
	}
}

func (s *OrderStore) Get(id string) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	order, exists := s.orders[id]
	return order, exists
}

func (s *OrderStore) Set(order Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.orders[order.ID] = order
}

func (s *OrderStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delete(s.orders, id)
}

func (s *OrderStore) GetAll() []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := make([]Order, 0, len(s.orders))
	for _, order := range s.orders {
		result = append(result, order)
	}
	return result
}

func (s *OrderStore) GetPendingOrders() []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var result []Order
	for _, order := range s.orders {
		if order.Status == "pending" {
			result = append(result, order)
		}
	}
	return result
}

var Store = NewOrderStore()
