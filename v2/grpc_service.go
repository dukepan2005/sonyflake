package sonyflake

import (
	"context"
)

var _ SonyflakeServiceServer = new(GrpcSonyflakeService)

// GrpcSonyflakeService implement SonyflakeServiceServer interface
type GrpcSonyflakeService struct {
	UnimplementedSonyflakeServiceServer
	IDCache chan int64
}

// NextID generate sequence
func (rpcS *GrpcSonyflakeService) NextID(ctx context.Context, req *SonyFlakeRequest) (*SonyFlakeResponse, error) {
	var (
		sequence int64
	)

	// generate ids, extract id from cache
	for idx := 0; idx < int(req.GetNum()); idx++ {
		sequence = <-rpcS.IDCache
	}

	return &SonyFlakeResponse{
		Sequence: sequence,
	}, nil
}
