package sonyflake

import (
	"context"
)

var _ SonyflakeServiceServer = new(GrpcSonyflakeService)

// GrpcSonyflakeService implement SonyflakeServiceServer interface
type GrpcSonyflakeService struct {
	UnimplementedSonyflakeServiceServer
	Sf *Sonyflake
}

// NextID generate sequence
func (rpcS *GrpcSonyflakeService) NextID(context.Context, *SonyFlakeRequest) (*SonyFlakeResponse, error) {
	id, err := rpcS.Sf.NextID()
	if err != nil {
		return nil, err
	}

	idParts := Decompose(id)

	return &SonyFlakeResponse{
		ID:        idParts["id"],
		Msb:       idParts["msg"],
		Time:      idParts["time"],
		Sequence:  idParts["sequence"],
		MachineID: idParts["machine-id"],
	}, nil
}
