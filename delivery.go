package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/garden-raccoon/delivery-pkg/models"
	proto "github.com/garden-raccoon/delivery-pkg/protocols/delivery-pkg"
)

type DeliveryPkgAPI interface {
	CreateOrUpdateDelivery(s *models.Delivery) error
	GetDeliveries() ([]*models.Delivery, error)
	DeleteDelivery(deliveryUuid uuid.UUID) error
	DeliveryById(deliveryUuid uuid.UUID) (*models.Delivery, error)
	HealthCheck() error
	// Close GRPC Api connection
	Close() error
}

// Api is profile-service GRPC Api
// structure with client Connection
type Api struct {
	addr    string
	timeout time.Duration
	*grpc.ClientConn
	mu sync.Mutex
	proto.DeliveryServiceClient
	grpc_health_v1.HealthClient
}

// New create new Battles Api instance
func New(addr string, timeOut time.Duration) (DeliveryPkgAPI, error) {
	api := &Api{timeout: timeOut * time.Second}

	if err := api.initConn(addr); err != nil {
		return nil, fmt.Errorf("create DeliveryApi:  %w", err)
	}
	api.HealthClient = grpc_health_v1.NewHealthClient(api.ClientConn)

	api.DeliveryServiceClient = proto.NewDeliveryServiceClient(api.ClientConn)
	return api, nil
}

func (api *Api) DeleteDelivery(deliveryUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	req := &proto.DeliveryDeleteReq{
		DeliveryUuid: deliveryUuid.Bytes(),
	}
	_, err := api.DeliveryServiceClient.DeleteDelivery(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteDelivery api request: %w", err)
	}
	return nil
}

func (api *Api) GetDeliveries() ([]*models.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	var resp *proto.Deliveries
	empty := new(proto.EmptyDelivery)
	resp, err := api.DeliveryServiceClient.GetDeliveries(ctx, empty)
	if err != nil {
		return nil, fmt.Errorf("GetDeliverys api request: %w", err)
	}

	deliverys, err := models.DeliveriesFromProto(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to GetDeliverys %w", err)
	}
	return deliverys, nil
}

func (api *Api) CreateOrUpdateDelivery(s *models.Delivery) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	deliverys, err := models.Proto(s)
	if err != nil {
		return fmt.Errorf("failed to CreateOrUpdateDelivery %w", err)
	}
	_, err = api.DeliveryServiceClient.CreateOrUpdateDelivery(ctx, deliverys)
	if err != nil {
		return fmt.Errorf("create delivery api request: %w", err)
	}
	return nil
}

// initConn initialize connection to Grpc servers
func (api *Api) initConn(addr string) (err error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	api.ClientConn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	return
}

func (api *Api) DeliveryById(deliveryUuid uuid.UUID) (*models.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	getReq := &proto.DeliveryGetReq{DeliveryUuid: deliveryUuid.Bytes()}
	resp, err := api.DeliveryServiceClient.DeliveryById(ctx, getReq)
	if err != nil {
		return nil, fmt.Errorf("DeliveryAPI DeliveryById request failed: %w", err)
	}

	delivery, err := models.DeliveryFromProto(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to DeliveryById %w", err)
	}
	return delivery, nil
}

func (api *Api) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	api.mu.Lock()
	defer api.mu.Unlock()

	resp, err := api.HealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "deliveryapi"})
	if err != nil {
		return fmt.Errorf("healthcheck error: %w", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("node is %s", errors.New("service is unhealthy"))
	}
	//api.status = service.StatusHealthy
	return nil
}
