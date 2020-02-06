package status

import (
	"bytes"
	"context"
	"github.com/codenotary/immudbrestproxy/pkg/api"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"io/ioutil"
)

const ROOT_FN = ".root"

type Rootprovider interface {
	GetRoot(ctx context.Context) (*api.Root, error)
	SetRoot(root *api.Root) error
}

type rootprovider struct {
	immuC api.ImmuServiceClient
}

func NewRootprovider(immuC api.ImmuServiceClient) Rootprovider {
	return &rootprovider{immuC}
}

func (r *rootprovider) GetRoot(ctx context.Context) (*api.Root, error) {
	root := new(api.Root)

	if buf, err := ioutil.ReadFile(ROOT_FN); err == nil {
		if err = root.XXX_Unmarshal(buf); err != nil {
			return nil, err
		}
		return root, nil
	}

	var protoReq empty.Empty
	var metadata runtime.ServerMetadata
	if root, err := r.immuC.CurrentRoot(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD)); err != nil {
		return nil, err
	} else {
		// Everytime we found a fresh root we persist it on ROOT_FN file
		if err := r.SetRoot(root); err != nil {
			return nil, err
		}
		return root, nil
	}
}

func (r *rootprovider) SetRoot(root *api.Root) error {
	var buf bytes.Buffer
	raw, err := root.XXX_Marshal(buf.Bytes(), true)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(ROOT_FN, raw, 0644)
	if err != nil {
		return err
	}
	return nil
}
