package api

import (
	"bytes"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"io/ioutil"
)

const ROOT_FN = ".root"

type Rootprovider interface {
	GetRoot(ctx context.Context) (*Root, error)
}

type rootprovider struct {
	immuC ImmuServiceClient
}

func NewRootprovider(immuC ImmuServiceClient) Rootprovider {
	return &rootprovider{immuC}
}

func (r *rootprovider) GetRoot(ctx context.Context) (*Root, error) {
	if root, err := GetCachedRoot(); err == nil {
		return root, nil
	}
	var protoReq empty.Empty
	var metadata runtime.ServerMetadata
	if root, err := r.immuC.CurrentRoot(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD)); err != nil {
		return nil, err
	} else {
		// Everytime we found a fresh root we persist it on ROOT_FN file
		if err := SetRoot(root); err != nil {
			return nil, err
		}
		return root, nil
	}
}

func SetRoot(root *Root) error {
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


func GetCachedRoot() (*Root, error) {
	root := new(Root)
	buf, err := ioutil.ReadFile(ROOT_FN)
	if err == nil {
		if err = root.XXX_Unmarshal(buf); err != nil {
			return nil, err
		}
		return root, nil
	}
	return nil, err
}
