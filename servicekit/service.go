package servicekit

import (
	"context"
	"sync"
)

// Service is a simple CRUD interface for engines.
type Service interface {
	RegisterEngine(ctx context.Context, e Engine) error
	GetRegisteredEngine(ctx context.Context, id string) (Engine, error)
}

// Engine CRUD data.
// ID should be globally unique.
type Engine struct {
	ID             string  `json:"id"`
	FactoryID      string  `json:"factory_id"`
	EngineConfig   string  `json:"engine_config"`
	EngineCapacity float32 `json:"engine_capacity"`
	FuelCapacity   float32 `json:"fuel_capacity"`
	FuelRange      float32 `json:"fuel_range"`
	EngineHP       float32 `json:"engine_hp"`
	EngineTorque   float32 `json:"engine_torque"`
}

type RegistrationService struct {
	mtx sync.RWMutex
	m   map[string]Engine
}

func NewRegistrationService() Service {
	return &RegistrationService{
		m: map[string]Engine{},
	}
}

func (s *RegistrationService) RegisterEngine(ctx context.Context, e Engine) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[e.ID]; ok {
		return ErrEnginePreExistInRegistry // POST = create, don't overwrite
	}
	s.m[e.ID] = e
	return nil
}

func (s *RegistrationService) GetRegisteredEngine(ctx context.Context, id string) (Engine, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	e, ok := s.m[id]
	if !ok {
		return Engine{}, ErrEngineNotExistInRegistry
	}
	return e, nil
}