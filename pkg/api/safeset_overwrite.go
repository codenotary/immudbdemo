// Copyright 2019, Oath Inc.
// Licensed under the terms of the Apache License 2.0. Please see LICENSE file in project root for terms.

package api

import (
	"context"
	"encoding/json"
	"github.com/codenotary/immudb/pkg/api/schema"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type ForwardSafeSetResp interface {
	forwardSafeSetResp(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error)
}

type forwardSafeSetResp struct {
}

func NewForwardSafeSetResp() ForwardSafeSetResp {
	return &forwardSafeSetResp{}
}

func (f forwardSafeSetResp) forwardSafeSetResp(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
	if req.Method == http.MethodPost && resp != nil && req.URL.Path == "/v1/immurestproxy/item/safe" {
		if p, ok := resp.(*Proof); ok {
			proof := schema.Proof{p.Leaf, p.Index, p.Root, p.At, p.InclusionPath, p.ConsistencyPath, p.XXX_NoUnkeyedLiteral, p.XXX_unrecognized, p.XXX_sizecache}
			root := new(schema.Root)

			rootc, _ := GetCachedRoot()
			// TODO FIX THIS. Try to use schema models
			root.Root = rootc.Root
			root.Index = rootc.Index

			w.Header().Set("Content-Type", "application/json")
			buf, err := marshaler.Marshal(resp)
			if err != nil {
				panic(err)
			}
			var m map[string]interface{}
			err = json.Unmarshal(buf, &m)
			if err != nil {
				panic(err)
			}
			/* remember to calc the leaf hash from key val with values that are coming from client and index from server.
			DO NOT USE leaf generated from server for security reasons. (maybe somebody can create a temper leaf)
			*/
			verified := proof.Verify(p.Leaf, *root)
			m["verified"] = verified
			newData, _ := json.Marshal(m)

			if verified {
				//saving a fresh root
				tocache := new(Root)
				tocache.Index = p.Index
				tocache.Root = p.Root
				SetRoot(tocache)
			}

			w.Write(newData)
			return
		}
	}
}

func init() {
	forwardSafeSetResp := NewForwardSafeSetResp()
	forward_ImmuService_SafeSet_0 = forwardSafeSetResp.forwardSafeSetResp
}
