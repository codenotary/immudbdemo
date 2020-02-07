// Copyright 2019, Oath Inc.
// Licensed under the terms of the Apache License 2.0. Please see LICENSE file in project root for terms.

package api

import (
	"context"
	"encoding/json"
	"github.com/codenotary/immudb/pkg/api/schema"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type ForwardSafeGetResp interface {
	forwardSafeGetResp(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error)
}

type forwardSafeGetResp struct {
}

func NewForwardSafeGetResp() ForwardSafeGetResp {
	return &forwardSafeGetResp{}
}

func (f forwardSafeGetResp) forwardSafeGetResp(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
	if req.Method == http.MethodPost && resp != nil && req.URL.Path == "/v1/immurestproxy/item/safe/get" {
		println(resp)
		if p, ok := resp.(*SafeItem); ok {
			proof := schema.Proof{p.Proof.Leaf, p.Proof.Index, p.Proof.Root, p.Proof.At, p.Proof.InclusionPath, p.Proof.ConsistencyPath, p.Proof.XXX_NoUnkeyedLiteral, p.Proof.XXX_unrecognized, p.Proof.XXX_sizecache}
			root := new(schema.Root)
			if buf, err := ioutil.ReadFile(".root"); err == nil {
				if err = root.XXX_Unmarshal(buf); err != nil {
					panic(err)
				}
			}

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
			verified := proof.Verify(p.Proof.Leaf, *root)
			m["verified"] = verified
			newData, _ := json.Marshal(m)

			if verified {
				//saving a fresh root
				tocache := new(Root)
				tocache.Index = p.Proof.Index
				tocache.Root = p.Proof.Root
				SetRoot(tocache)
			}

			w.Write(newData)
			return
		}
	}
}

func init() {
	forwardSafeGetResp := NewForwardSafeGetResp()
	forward_ImmuService_SafeGet_0 = forwardSafeGetResp.forwardSafeGetResp
}
