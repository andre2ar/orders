package service

import (
	"context"

	"github.com/andre2ar/orders/internal/infra/grpc/pb"
	"github.com/andre2ar/orders/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, b *pb.Blank) (*pb.OrdersList, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var listOfOrders []*pb.Order
	for _, order := range orders {
		listOfOrders = append(listOfOrders, &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}

	return &pb.OrdersList{Orders: listOfOrders}, nil
}
