package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/GrudTrigger/go-grpc-graphql-microservice/account"
	"github.com/GrudTrigger/go-grpc-graphql-microservice/catalog"
	"github.com/GrudTrigger/go-grpc-graphql-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service: s,
		accountClient,
		catalogClient,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
} 

func(s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest)(*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil{
		log.Println("Error getting account: ", err)
		return nil, errors.New("account not found")
	}

	productsIDs := []string{}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productsIDs, "")
	if err != nil {
		log.Println("Error getting products: ", err)
		return nil, errors.New("products not found")
	}
	products := []OrderedProduct{}
	for _, p := range orderedProducts{
		product := orderedProducts{
			ID: p.ID,
			Quantity: 0,
			Price: p.Price,
			Name: p.Name,
			Description: p.Description,
		}
		for _, rp := range r.Products{
			if rp.ProductId == p.ID{
				product.Quantity = rp.Quantity
				break
			}
		}
		if product.Quantity != 0 {
			products = append(products, product)
		}
	}
	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Println("Error posting order:", err)
		return nil, errors.New("could not post order")
	}

	orderProto := &pb.Order{
		Id: order.ID,
		AccountId: order.AccountID,
		TotalPrice: order.TotalPrice,
		Products: []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products{
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id: p.ID,
			Name: p.Name,
			Description: p.Description,
			Price: p.Price,
			Quantity: p.Quantity,
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func(s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrderForAccountRequest) (*pb.GetOrdersForAccountResponse, error){
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIDMap := map[string]bool{}
	for _, o := range accaccountOrders{
		for _, p := range o.Products{
			productIDMap[p.ID] = true
		}
	}
	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting account products: ", err)
		return nil, err
	}
	orders := []*pb.Order{}
	range accaccountOrders{
		
	}
}