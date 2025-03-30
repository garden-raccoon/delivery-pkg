package models

import (
	"fmt"
	proto "github.com/garden-raccoon/delivery-pkg/protocols/delivery-pkg"
	"github.com/gofrs/uuid"
)

type Delivery struct {
	DeliveryUuid uuid.UUID `json:"delivery_uuid"`
}

// Proto is
func Proto(deli *Delivery) (*proto.Delivery, error) {
	d := &proto.Delivery{
		DeliveryUuid: deli.DeliveryUuid.Bytes(),
	}
	return d, nil
}

func DeliveryFromProto(pb *proto.Delivery) (*Delivery, error) {
	meal := &Delivery{
		DeliveryUuid: uuid.FromBytesOrNil(pb.DeliveryUuid),
	}
	return meal, nil
}

func DeliveriesToProto(delivery []*Delivery) (*proto.Deliveries, error) {
	pb := &proto.Deliveries{}
	for _, b := range delivery {
		meal, err := Proto(b)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		pb.Deliveries = append(pb.Deliveries, meal)
	}
	return pb, nil
}

func DeliveriesFromProto(pb *proto.Deliveries) ([]*Delivery, error) {
	var delivery []*Delivery
	for _, b := range pb.Deliveries {
		meal, err := DeliveryFromProto(b)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		delivery = append(delivery, meal)
	}
	return delivery, nil
}
