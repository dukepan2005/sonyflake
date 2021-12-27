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
func (rpcS *GrpcSonyflakeService) NextID(ctx context.Context, req *SonyFlakeRequest) (*SonyFlakeResponse, error) {
	var (
		sequence uint64
		err      error
	)

	for idx := 0; idx < int(req.GetNum()); idx++ {
		sequence, err = rpcS.Sf.NextID()
		if err != nil {
			return nil, err
		}
	}

	idParts := Decompose(sequence)

	return &SonyFlakeResponse{
		Id:        idParts["id"],
		Msb:       idParts["msb"],
		Time:      idParts["time"],
		Sequence:  idParts["sequence"],
		MachineID: idParts["machine-id"],
	}, nil
}
