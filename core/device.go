package core

import (
	"context"

	"github.com/altercpre/aparok-sync/app"

	"github.com/altercpre/aparok-core/core/device"
	"github.com/altercpre/aparok-core/pb"
)

func (mw *Middleware) DeviceNetworkStateSet(cctx context.Context, req *pb.RpcDeviceNetworkStateSetRequest) *pb.RpcDeviceNetworkStateSetResponse {
	response := func(code pb.RpcDeviceNetworkStateSetResponseErrorCode, err error) *pb.RpcDeviceNetworkStateSetResponse {
		m := &pb.RpcDeviceNetworkStateSetResponse{Error: &pb.RpcDeviceNetworkStateSetResponseError{Code: code}}
		if err != nil {
			m.Error.Description = getErrorDescription(err)
		}
		return m
	}
	mwApp := mw.GetApp()
	if mwApp == nil {
		return response(pb.RpcDeviceNetworkStateSetResponseError_INTERNAL_ERROR, ErrNotLoggedIn)
	}
	app.MustComponent[device.NetworkState](mwApp).SetNetworkState(req.DeviceNetworkType, req.NetworkId)
	return response(pb.RpcDeviceNetworkStateSetResponseError_NULL, nil)
}
