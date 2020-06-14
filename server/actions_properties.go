// Copyright 2020 Google LLC. All Rights Reserved.

package server

import (
	"context"
	"log"

	"apigov.dev/registry/models"
	rpc "apigov.dev/registry/rpc"
	"cloud.google.com/go/datastore"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *RegistryServer) CreateProperty(ctx context.Context, request *rpc.CreatePropertyRequest) (*rpc.Property, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, internalError(err)
	}
	defer client.Close()
	property, err := models.NewPropertyFromParentAndPropertyID(request.GetParent(), request.GetPropertyId())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	err = property.Update(request.GetProperty())
	k := &datastore.Key{Kind: models.PropertyEntityName, Name: property.ResourceName()}
	// fail if property already exists
	var existingProperty models.Property
	err = client.Get(ctx, k, &existingProperty)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, property.ResourceName()+" already exists")
	}
	property.CreateTime = property.UpdateTime
	k, err = client.Put(ctx, k, property)
	if err != nil {
		return nil, internalError(err)
	}
	return property.Message()
}

func (s *RegistryServer) DeleteProperty(ctx context.Context, request *rpc.DeletePropertyRequest) (*empty.Empty, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, internalError(err)
	}
	defer client.Close()
	// Validate name and create dummy property (we just need the ID fields).
	_, err = models.NewPropertyFromResourceName(request.GetName())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	// Delete the property.
	k := &datastore.Key{Kind: models.PropertyEntityName, Name: request.GetName()}
	err = client.Delete(ctx, k)
	return &empty.Empty{}, internalError(err)
}

func (s *RegistryServer) GetProperty(ctx context.Context, request *rpc.GetPropertyRequest) (*rpc.Property, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, internalError(err)
	}
	defer client.Close()
	property, err := models.NewPropertyFromResourceName(request.GetName())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	log.Printf("looking for %s", property.ResourceName())
	k := &datastore.Key{Kind: models.PropertyEntityName, Name: property.ResourceName()}
	err = client.Get(ctx, k, property)
	if err == datastore.ErrNoSuchEntity {
		return nil, status.Error(codes.NotFound, "not found")
	} else if err != nil {
		return nil, internalError(err)
	}
	return property.Message()
}

func (s *RegistryServer) ListProperties(ctx context.Context, req *rpc.ListPropertiesRequest) (*rpc.ListPropertiesResponse, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, internalError(err)
	}
	defer client.Close()
	q := datastore.NewQuery(models.PropertyEntityName)
	q, err = queryApplyCursor(q, req.GetPageToken())
	if err != nil {
		return nil, internalError(err)
	}
	p, err := models.NewPropertyFromParentAndPropertyID(req.GetParent(), "-")
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	if p.ProjectID != "-" && p.ProjectID != "-" {
		q = q.Filter("ProjectID =", p.ProjectID)
	}
	if p.ProductID != "-" && p.ProductID != "-" {
		q = q.Filter("ProductID =", p.ProductID)
	}
	if p.VersionID != "-" && p.VersionID != "-" {
		q = q.Filter("VersionID =", p.VersionID)
	}
	if p.SpecID != "-" && p.SpecID != "-" {
		q = q.Filter("SpecID =", p.SpecID)
	}
	prg, err := createFilterOperator(req.GetFilter(),
		[]filterArg{
			{"project_id", filterArgTypeString},
			{"product_id", filterArgTypeString},
			{"version_id", filterArgTypeString},
			{"spec_id", filterArgTypeString},
			{"property_id", filterArgTypeString},
		})
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	var propertyMessages []*rpc.Property
	var property models.Property
	it := client.Run(ctx, q.Distinct())
	pageSize := boundPageSize(req.GetPageSize())
	for _, err = it.Next(&property); err == nil; _, err = it.Next(&property) {
		if prg != nil {
			out, _, err := prg.Eval(map[string]interface{}{
				"project_id":  property.ProjectID,
				"product_id":  property.ProductID,
				"version_id":  property.VersionID,
				"spec_id":     property.SpecID,
				"property_id": property.PropertyID,
			})
			if err != nil {
				return nil, invalidArgumentError(err)
			}
			if !out.Value().(bool) {
				continue
			}
		}
		propertyMessage, _ := property.Message()
		propertyMessages = append(propertyMessages, propertyMessage)
		if len(propertyMessages) == pageSize {
			break
		}
	}
	if err != nil && err != iterator.Done {
		return nil, internalError(err)
	}
	responses := &rpc.ListPropertiesResponse{
		Properties: propertyMessages,
	}
	responses.NextPageToken, err = iteratorGetCursor(it, len(propertyMessages))
	if err != nil {
		return nil, internalError(err)
	}
	return responses, nil
}

func (s *RegistryServer) UpdateProperty(ctx context.Context, request *rpc.UpdatePropertyRequest) (*rpc.Property, error) {
	client, err := s.newDataStoreClient(ctx)
	if err != nil {
		return nil, internalError(err)
	}
	defer client.Close()
	property, err := models.NewPropertyFromResourceName(request.GetProperty().GetName())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	k := &datastore.Key{Kind: models.PropertyEntityName, Name: request.GetProperty().GetName()}
	err = client.Get(ctx, k, property)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	err = property.Update(request.GetProperty())
	if err != nil {
		return nil, internalError(err)
	}
	k, err = client.Put(ctx, k, property)
	if err != nil {
		return nil, internalError(err)
	}
	return property.Message()
}
