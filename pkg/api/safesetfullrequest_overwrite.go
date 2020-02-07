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

type RequestSafeSetCustom interface {
	request_ImmuService_SafeSet_Custom(context.Context, runtime.Marshaler, ImmuServiceClient, *http.Request, map[string]string) (proto.Message, runtime.ServerMetadata, error)
}

type requestSafeSetCustom struct {
}

func NewRequestSafeSetCustom() RequestSafeSetCustom {
	return &requestSafeSetCustom{}
}

func (r requestSafeSetCustom) request_ImmuService_SafeSet_Custom(ctx context.Context, marshaler runtime.Marshaler, client ImmuServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SafeSetOptions
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
	log.Print(fmt.Sprintf("Inserting:\n index: %d \nkey: %s\n val: %s", protoReq.RootIndex.Index, string(protoReq.Kv.Key), string(protoReq.Kv.Value)))

	msg, err := client.SafeSet(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func init() {
	r := NewRequestSafeSetCustom()
	overwrite_request_ImmuService_SafeSet_0 = r.request_ImmuService_SafeSet_Custom
}
