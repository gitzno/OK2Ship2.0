package registry

import "export_module_service/internal/domain"

type UseCaseRegistry struct {
	Handlers map[string]domain.MainUseCase
}

func NewUseCaseRegistry(peel domain.MainUseCase) *UseCaseRegistry {
	return &UseCaseRegistry{
		Handlers: map[string]domain.MainUseCase{
			"peel_test": peel,
			// "tension_test": tension,
		},
	}
}
