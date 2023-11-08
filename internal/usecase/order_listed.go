package usecase

import (
	"github.com/kenesparta/fullcycle-clean-architecture/internal/dto"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/entity"
	"github.com/kenesparta/fullcycle-clean-architecture/pkg/events"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrderUseCase) Execute() ([]dto.OrderOutput, error) {
	orders, listErr := c.OrderRepository.List()
	if listErr != nil {
		return nil, listErr
	}

	var orderOutputList []dto.OrderOutput
	for _, order := range orders {
		orderOutputList = append(orderOutputList, dto.OrderOutput{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return orderOutputList, nil
}
