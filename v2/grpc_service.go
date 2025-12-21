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
		sequence int64
		err      error
	)

	for idx := 0; idx < int(req.GetNum()); idx++ {
		sequence, err = rpcS.Sf.NextID()
		if err != nil {
			return nil, err
		}
	}

	return &SonyFlakeResponse{
		Sequence: sequence,
	}, nil
}
