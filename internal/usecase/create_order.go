package usecase

import (
	"github.com/kenesparta/fullcycle-clean-architecture/internal/dto"
	"github.com/kenesparta/fullcycle-clean-architecture/internal/entity"
	"github.com/kenesparta/fullcycle-clean-architecture/pkg/events"
)

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input dto.OrderInput) (dto.OrderOutput, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}

	if err := order.CalculateFinalPrice(); err != nil {
		return dto.OrderOutput{}, err
	}

	if err := c.OrderRepository.Save(&order); err != nil {
		return dto.OrderOutput{}, err
	}

	orderOutput := dto.OrderOutput{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	c.OrderCreated.SetPayload(orderOutput)
	dispErr := c.EventDispatcher.Dispatch(c.OrderCreated)
	if dispErr != nil {
		return dto.OrderOutput{}, dispErr
	}

	return orderOutput, nil
}
