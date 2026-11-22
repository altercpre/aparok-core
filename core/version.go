package core

import (
	"context"

	"github.com/altercpre/aparok-core/pb"
	"github.com/altercpre/aparok-core/util/vcs"
)

func (mw *Middleware) AppGetVersion(cctx context.Context, req *pb.RpcAppGetVersionRequest) *pb.RpcAppGetVersionResponse {
	response := func(version, details string, code pb.RpcAppGetVersionResponseErrorCode, err error) *pb.RpcAppGetVersionResponse {
		m := &pb.RpcAppGetVersionResponse{Version: version, Details: details, Error: &pb.RpcAppGetVersionResponseError{Code: code}}
		if err != nil {
			m.Error.Description = getErrorDescription(err)
		}

		return m
	}

	info := vcs.GetVCSInfo()
	return response(info.Version(), info.Description(), pb.RpcAppGetVersionResponseError_NULL, nil)
}
