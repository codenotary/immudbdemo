// Copyright 2019, Oath Inc.
// Licensed under the terms of the Apache License 2.0. Please see LICENSE file in project root for terms.

package api

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
)

type RequestSafeGetCustom interface {
	request_ImmuService_SafeGet_Custom(context.Context, runtime.Marshaler, ImmuServiceClient, *http.Request, map[string]string) (proto.Message, runtime.ServerMetadata, error)
}

type requestSafeGetCustom struct {
}

func NewRequestSafeGetCustom() RequestSafeGetCustom {
	return &requestSafeGetCustom{}
}
func (r requestSafeGetCustom) request_ImmuService_SafeGet_Custom(ctx context.Context, marshaler runtime.Marshaler, client ImmuServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SafeGetOptions
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	root, err := GetCachedRoot()

	ri := new(Index)
	ri.Index = root.Index
	protoReq.RootIndex = ri
	log.Print(fmt.Sprintf("Retrieving:\n index: %d \nkey: %s", protoReq.RootIndex.Index, string(protoReq.Key.Key)))

	msg, err := client.SafeGet(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func init() {
	r := NewRequestSafeGetCustom()
	overwrite_request_ImmuService_SafeGet_0 = r.request_ImmuService_SafeGet_Custom
}
