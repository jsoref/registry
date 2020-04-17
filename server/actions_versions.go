// Copyright 2020 Google LLC. All Rights Reserved.

package server

import (
	"context"
	"log"

	"apigov.dev/flame/models"
	rpc "apigov.dev/flame/rpc"
	"cloud.google.com/go/datastore"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *FlameServer) CreateVersion(ctx context.Context, request *rpc.CreateVersionRequest) (*rpc.Version, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	version, err := models.NewVersionFromParentAndVersionID(request.GetParent(), request.GetVersionId())
	if err != nil {
		return nil, err
	}
	k := &datastore.Key{Kind: models.VersionEntityName, Name: version.ResourceName()}
	// fail if version already exists
	var existingVersion models.Version
	err = client.Get(ctx, k, &existingVersion)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, version.ResourceName()+" already exists")
	}
	err = version.Update(request.GetVersion())
	version.CreateTime = version.UpdateTime
	k, err = client.Put(ctx, k, version)
	if err != nil {
		return nil, err
	}
	return version.Message()
}

func (s *FlameServer) DeleteVersion(ctx context.Context, request *rpc.DeleteVersionRequest) (*empty.Empty, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	// Validate name and create dummy version (we just need the ID fields).
	version, err := models.NewVersionFromResourceName(request.GetName())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	// Delete children first and then delete the version.
	version.DeleteChildren(ctx, client)
	k := &datastore.Key{Kind: models.VersionEntityName, Name: request.GetName()}
	err = client.Delete(ctx, k)
	return &empty.Empty{}, err
}

func (s *FlameServer) GetVersion(ctx context.Context, request *rpc.GetVersionRequest) (*rpc.Version, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	version, err := models.NewVersionFromResourceName(request.GetName())
	if err != nil {
		return nil, err
	}
	k := &datastore.Key{Kind: models.VersionEntityName, Name: version.ResourceName()}
	err = client.Get(ctx, k, version)
	if err == datastore.ErrNoSuchEntity {
		return nil, status.Error(codes.NotFound, "not found")
	} else if err != nil {
		return nil, internalError(err)
	}
	return version.Message()
}

func (s *FlameServer) ListVersions(ctx context.Context, req *rpc.ListVersionsRequest) (*rpc.ListVersionsResponse, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	log.Printf("list versions")
	q := datastore.NewQuery(models.VersionEntityName)
	q, err = queryApplyCursor(q, req.GetPageToken())
	if err != nil {
		return nil, internalError(err)
	}
	m, err := models.ParseParentProduct(req.GetParent())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	if m[1] != "-" {
		q = q.Filter("ProjectID =", m[1])
	}
	if m[2] != "-" {
		q = q.Filter("ProductID =", m[2])
	}
	prg, err := createFilterOperator(req.GetFilter(),
		[]filterArg{
			{"version_id", filterArgTypeString},
			{"state", filterArgTypeString},
		})
	if err != nil {
		return nil, internalError(err)
	}
	var versionMessages []*rpc.Version
	var version models.Version
	log.Printf("running %+v", q)
	it := client.Run(ctx, q.Distinct())
	pageSize := boundPageSize(req.GetPageSize())
	for _, err = it.Next(&version); err == nil; _, err = it.Next(&version) {
		if prg != nil {
			out, _, err := prg.Eval(map[string]interface{}{
				"version_id": version.VersionID,
				"state":      version.State,
			})
			if err != nil {
				return nil, invalidArgumentError(err)
			}
			if !out.Value().(bool) {
				continue
			}
		}
		versionMessage, _ := version.Message()
		versionMessages = append(versionMessages, versionMessage)
		if len(versionMessages) == pageSize {
			break
		}
	}
	if err != nil && err != iterator.Done {
		return nil, internalError(err)
	}
	responses := &rpc.ListVersionsResponse{
		Versions: versionMessages,
	}
	responses.NextPageToken, err = iteratorGetCursor(it, len(versionMessages))
	if err != nil {
		return nil, internalError(err)
	}
	return responses, nil
}

func (s *FlameServer) UpdateVersion(ctx context.Context, request *rpc.UpdateVersionRequest) (*rpc.Version, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	version, err := models.NewVersionFromResourceName(request.GetVersion().GetName())
	if err != nil {
		return nil, err
	}
	k := &datastore.Key{Kind: models.VersionEntityName, Name: version.ResourceName()}
	err = client.Get(ctx, k, &version)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	err = version.Update(request.GetVersion())
	if err != nil {
		return nil, err
	}
	k, err = client.Put(ctx, k, version)
	if err != nil {
		return nil, err
	}
	return version.Message()
}
