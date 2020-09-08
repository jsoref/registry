// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/models"
	"github.com/apigee/registry/server/names"
	storage "github.com/apigee/registry/server/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateSpec handles the corresponding API request.
func (s *RegistryServer) CreateSpec(ctx context.Context, request *rpc.CreateSpecRequest) (*rpc.Spec, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	spec, err := models.NewSpecFromParentAndSpecID(request.GetParent(), request.GetSpecId())
	if err != nil {
		return nil, internalError(err)
	}
	// fail if spec already exists
	q := client.NewQuery(models.SpecEntityName)
	q = q.Require("ProjectID", spec.ProjectID)
	q = q.Require("ApiID", spec.ApiID)
	q = q.Require("VersionID", spec.VersionID)
	q = q.Require("SpecID", spec.SpecID)
	it := client.Run(ctx, q)
	var existingSpec models.Spec
	existingKey, err := it.Next(&existingSpec)
	if existingKey != nil {
		return nil, status.Error(codes.AlreadyExists, spec.ResourceName()+" already exists")
	}
	// save the spec under its full resource@revision name
	err = spec.Update(request.GetSpec())
	if err != nil {
		return nil, internalError(err)
	}
	spec.CreateTime = spec.UpdateTime
	// the first revision of the spec that we save is also the current one
	spec.Currency = models.IsCurrent
	k := client.NewKey(models.SpecEntityName, spec.ResourceNameWithRevision())
	k, err = client.Put(ctx, k, spec)
	if err != nil {
		log.Printf("save spec error %+v", err)
		return nil, internalError(err)
	}
	// save a blob with the spec contents
	blob := models.NewBlob(
		spec,
		request.GetSpec().GetContents())
	k2 := client.NewKey(models.BlobEntityName, spec.ResourceNameWithRevision())
	_, err = client.Put(ctx,
		k2,
		blob)
	if err != nil {
		log.Printf("save blob error %+v", err)
		return nil, internalError(err)
	}
	response, nil := spec.Message(nil, "")
	s.notify(rpc.Notification_CREATED, spec.ResourceNameWithRevision())
	return response, nil
}

// DeleteSpec handles the corresponding API request.
func (s *RegistryServer) DeleteSpec(ctx context.Context, request *rpc.DeleteSpecRequest) (*empty.Empty, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	// Validate name and create dummy spec (we just need the ID fields).
	spec, err := models.NewSpecFromResourceName(request.GetName())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	if spec.RevisionID != "" {
		return nil, invalidArgumentError(errors.New("specific revisions should be deleted with DeleteSpecRevision"))
	}
	// Delete all revisions of the spec.
	q := client.NewQuery(models.SpecEntityName)
	q = q.Require("ProjectID", spec.ProjectID)
	q = q.Require("ApiID", spec.ApiID)
	q = q.Require("VersionID", spec.VersionID)
	q = q.Require("SpecID", spec.SpecID)
	err = client.DeleteAllMatches(ctx, q)
	// Delete all blobs associated with the spec.
	q = client.NewQuery(models.BlobEntityName)
	q = q.Require("ProjectID", spec.ProjectID)
	q = q.Require("ApiID", spec.ApiID)
	q = q.Require("VersionID", spec.VersionID)
	q = q.Require("SpecID", spec.SpecID)
	err = client.DeleteAllMatches(ctx, q)
	s.notify(rpc.Notification_DELETED, request.GetName())
	return &empty.Empty{}, err
}

// GetSpec handles the corresponding API request.
func (s *RegistryServer) GetSpec(ctx context.Context, request *rpc.GetSpecRequest) (*rpc.Spec, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	spec, userSpecifiedRevision, err := fetchSpec(ctx, client, request.GetName())
	if err != nil {
		if client.IsNotFound(err) {
			return nil, notFoundError(err)
		}
		return nil, internalError(err)
	}
	var blob *models.Blob
	if request.GetView() == rpc.SpecView_FULL {
		blob, _ = fetchBlobForSpec(ctx, client, spec)
	}
	return spec.Message(blob, userSpecifiedRevision)
}

// ListSpecs handles the corresponding API request.
func (s *RegistryServer) ListSpecs(ctx context.Context, req *rpc.ListSpecsRequest) (*rpc.ListSpecsResponse, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	q := client.NewQuery(models.SpecEntityName)
	q, err = q.ApplyCursor(req.GetPageToken())
	if err != nil {
		return nil, internalError(err)
	}
	m, err := names.ParseParentVersion(req.GetParent())
	if err != nil {
		return nil, invalidArgumentError(err)
	}
	if m[1] != "-" {
		q = q.Require("ProjectID", m[1])
	}
	if m[2] != "-" {
		q = q.Require("ApiID", m[2])
	}
	if m[3] != "-" {
		q = q.Require("VersionID", m[3])
	}
	q = q.Require("Currency", models.IsCurrent)
	prg, err := createFilterOperator(req.GetFilter(),
		[]filterArg{
			{"project_id", filterArgTypeString},
			{"api_id", filterArgTypeString},
			{"version_id", filterArgTypeString},
			{"spec_id", filterArgTypeString},
			{"filename", filterArgTypeString},
			{"description", filterArgTypeString},
			{"style", filterArgTypeString},
		})
	if err != nil {
		return nil, internalError(err)
	}
	var specMessages []*rpc.Spec
	var spec models.Spec
	it := client.Run(ctx, q)
	pageSize := boundPageSize(req.GetPageSize())
	for _, err := it.Next(&spec); err == nil; _, err = it.Next(&spec) {
		if prg != nil {
			out, _, err := prg.Eval(map[string]interface{}{
				"project_id":  spec.ProjectID,
				"api_id":      spec.ApiID,
				"version_id":  spec.VersionID,
				"spec_id":     spec.SpecID,
				"filename":    spec.FileName,
				"description": spec.Description,
				"style":       spec.Style,
			})
			if err != nil {
				return nil, invalidArgumentError(err)
			}
			if !out.Value().(bool) {
				continue
			}
		}
		var blob *models.Blob
		if req.GetView() == rpc.SpecView_FULL {
			blob, _ = fetchBlobForSpec(ctx, client, &spec)
		}
		specMessage, _ := spec.Message(blob, "")
		specMessages = append(specMessages, specMessage)
		if len(specMessages) == pageSize {
			break
		}
	}
	if err != nil && err != iterator.Done {
		return nil, internalError(err)
	}
	responses := &rpc.ListSpecsResponse{
		Specs: specMessages,
	}
	responses.NextPageToken, err = it.GetCursor(len(specMessages))
	if err != nil {
		return nil, internalError(err)
	}
	return responses, nil
}

// UpdateSpec handles the corresponding API request.
func (s *RegistryServer) UpdateSpec(ctx context.Context, request *rpc.UpdateSpecRequest) (*rpc.Spec, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	spec, userSpecifiedRevision, err := fetchSpec(ctx, client, request.GetSpec().GetName())
	if err != nil {
		return nil, internalError(err)
	}
	if userSpecifiedRevision != "" {
		return nil, invalidArgumentError(errors.New("updates to specific revisions are unsupported"))
	}
	oldRevisionID := spec.RevisionID
	err = spec.Update(request.GetSpec())
	if err != nil {
		return nil, internalError(err)
	}
	newRevisionID := spec.RevisionID
	// if the revision changed, get the previously-current revision and mark it as non-current
	if oldRevisionID != newRevisionID {
		k := client.NewKey(models.SpecEntityName, spec.ResourceNameWithSpecifiedRevision(oldRevisionID))
		currentRevision := &models.Spec{}
		client.Get(ctx, k, currentRevision)
		currentRevision.Currency = models.NotCurrent
		_, err = client.Put(ctx, k, currentRevision)
		if err != nil {
			return nil, internalError(err)
		}
		spec.Currency = models.IsCurrent
	}
	k := client.NewKey(models.SpecEntityName, spec.ResourceNameWithRevision())
	spec.Key = spec.ResourceNameWithRevision()
	k, err = client.Put(ctx, k, spec)
	if err != nil {
		return nil, internalError(err)
	}
	// save a blob with the spec contents (but only if the contents were updated)
	if request.GetSpec().GetContents() != nil {
		blob := models.NewBlob(
			spec,
			request.GetSpec().GetContents())
		_, err = client.Put(ctx,
			client.NewKey(models.BlobEntityName, spec.ResourceNameWithRevision()),
			blob)
		if err != nil {
			return nil, internalError(err)
		}
	}
	s.notify(rpc.Notification_UPDATED, spec.ResourceNameWithRevision())
	return spec.Message(nil, "")
}

// ListSpecRevisions handles the corresponding API request.
func (s *RegistryServer) ListSpecRevisions(ctx context.Context, req *rpc.ListSpecRevisionsRequest) (*rpc.ListSpecRevisionsResponse, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	targetSpec, err := models.NewSpecFromResourceName(req.GetName())
	if err != nil {
		return nil, internalError(err)
	}
	q := client.NewQuery(models.SpecEntityName)
	q, err = q.ApplyCursor(req.GetPageToken())
	if err != nil {
		return nil, internalError(err)
	}
	q = q.Require("ProjectID", targetSpec.ProjectID)
	q = q.Require("ApiID", targetSpec.ApiID)
	q = q.Require("VersionID", targetSpec.VersionID)
	q = q.Require("SpecID", targetSpec.SpecID)
	q = q.Order("-CreateTime")

	var specMessages []*rpc.Spec
	responses := &rpc.ListSpecRevisionsResponse{}
	if s.weTrustTheSort {
		var spec models.Spec
		it := client.Run(ctx, q)
		pageSize := boundPageSize(req.GetPageSize())
		for _, err := it.Next(&spec); err == nil; _, err = it.Next(&spec) {
			specMessage, _ := spec.Message(nil, spec.RevisionID)
			specMessages = append(specMessages, specMessage)
			if len(specMessages) == pageSize {
				break
			}
		}
		if err != nil && err != iterator.Done {
			return nil, internalError(err)
		}
		responses.NextPageToken, err = it.GetCursor(len(specMessages))
		if err != nil {
			return nil, internalError(err)
		}
	} else {
		specs := make([]*models.Spec, 0)
		it := client.Run(ctx, q)
		for {
			spec := &models.Spec{}
			_, err := it.Next(spec)
			if err != nil {
				break
			}
			specs = append(specs, spec)
		}
		if err != nil && err != iterator.Done {
			return nil, internalError(err)
		}
		sort.Slice(specs, func(i, j int) bool {
			return specs[i].CreateTime.After(specs[j].CreateTime)
		})
		for _, spec := range specs {
			specMessage, _ := spec.Message(nil, spec.RevisionID)
			specMessages = append(specMessages, specMessage)
		}
		responses.NextPageToken = ""
		err = nil
	}
	responses.Specs = specMessages
	return responses, nil
}

// DeleteSpecRevision handles the corresponding API request.
func (s *RegistryServer) DeleteSpecRevision(ctx context.Context, request *rpc.DeleteSpecRevisionRequest) (*empty.Empty, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	// Delete the spec revision.
	// First, get the revision to delete.
	spec, _, err := fetchSpec(ctx, client, request.GetName())
	if err != nil {
		return nil, internalError(err)
	}
	k := client.NewKey(models.SpecEntityName, spec.ResourceNameWithRevision())
	// If the one we will delete is the current revision, we need to designate a new current revision.
	if spec.Currency == models.IsCurrent {
		// get the most recent non-current revision and make it current
		newKey, newCurrentRevision, err := s.fetchMostRecentNonCurrentRevisionOfSpec(ctx, client, request.GetName())
		if err != nil {
			log.Printf("error %+v", err)
		}
		if err == nil && newCurrentRevision != nil {
			newCurrentRevision.Currency = models.IsCurrent
			client.Put(ctx, newKey, newCurrentRevision)
		}
	}
	err = client.Delete(ctx, k)
	// Delete the blob associated with the spec
	k2 := client.NewKey(models.BlobEntityName, spec.ResourceNameWithRevision())
	err = client.Delete(ctx, k2)
	s.notify(rpc.Notification_DELETED, spec.ResourceNameWithRevision())
	return &empty.Empty{}, err
}

// TagSpecRevision handles the corresponding API request.
func (s *RegistryServer) TagSpecRevision(ctx context.Context, request *rpc.TagSpecRevisionRequest) (*rpc.Spec, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	spec, userSpecifiedRevision, err := fetchSpec(ctx, client, request.GetName())
	if err != nil {
		return nil, internalError(err)
	}
	if userSpecifiedRevision == "" {
		log.Printf("we might not want to support tagging specs with unspecified revisions")
	}
	if request.GetTag() == "" {
		return nil, invalidArgumentError(errors.New("tags cannot be empty"))
	}
	// save the tag
	now := time.Now()
	tag := &models.SpecRevisionTag{
		ProjectID:  spec.ProjectID,
		ApiID:      spec.ApiID,
		VersionID:  spec.VersionID,
		SpecID:     spec.SpecID,
		RevisionID: spec.RevisionID,
		Tag:        request.GetTag(),
		CreateTime: now,
		UpdateTime: now,
	}
	k := client.NewKey(models.SpecRevisionTagEntityName, tag.ResourceNameWithTag())
	k, err = client.Put(ctx, k, tag)
	// send a notification that the tagged spec has been updated
	s.notify(rpc.Notification_UPDATED, spec.ResourceNameWithSpecifiedRevision(request.GetTag()))
	// return the spec using the tag for its name
	return spec.Message(nil, request.GetTag())
}

// ListSpecRevisionTags handles the corresponding API request.
func (s *RegistryServer) ListSpecRevisionTags(ctx context.Context, req *rpc.ListSpecRevisionTagsRequest) (*rpc.ListSpecRevisionTagsResponse, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	targetSpec, err := models.NewSpecFromResourceName(req.GetName())
	if err != nil {
		return nil, internalError(err)
	}
	q := client.NewQuery(models.SpecRevisionTagEntityName)
	q, err = q.ApplyCursor(req.GetPageToken())
	if err != nil {
		return nil, internalError(err)
	}
	q = q.Require("ProjectID", targetSpec.ProjectID)
	q = q.Require("ApiID", targetSpec.ApiID)
	q = q.Require("VersionID", targetSpec.VersionID)
	q = q.Require("SpecID", targetSpec.SpecID)
	var tagMessages []*rpc.SpecRevisionTag
	tag := models.SpecRevisionTag{}
	it := client.Run(ctx, q)
	pageSize := boundPageSize(req.GetPageSize())
	for _, err := it.Next(&tag); err == nil; _, err = it.Next(&tag) {
		tagMessage, _ := tag.Message()
		tagMessages = append(tagMessages, tagMessage)
		if len(tagMessages) == pageSize {
			break
		}
	}
	if err != nil && err != iterator.Done {
		return nil, internalError(err)
	}
	responses := &rpc.ListSpecRevisionTagsResponse{
		Tags: tagMessages,
	}
	responses.NextPageToken, err = it.GetCursor(len(tagMessages))
	if err != nil {
		return nil, internalError(err)
	}
	return responses, nil
}

// RollbackSpec handles the corresponding API request.
func (s *RegistryServer) RollbackSpec(ctx context.Context, request *rpc.RollbackSpecRequest) (*rpc.Spec, error) {
	client, err := s.getStorageClient(ctx)
	if err != nil {
		return nil, unavailableError(err)
	}
	defer s.releaseStorageClient(client)
	specNameWithRevision := request.GetName() + "@" + request.GetRevisionId()
	spec, userSpecifiedRevision, err := fetchSpec(ctx, client, specNameWithRevision)
	if err != nil {
		// TODO: this should return NotFound if the revision was not found.
		return nil, notFoundError(err)
	}
	if userSpecifiedRevision == "" {
		return nil, invalidArgumentError(errors.New("rollbacks require a specified revision"))
	}
	// The previous current revision needs to be marked non-current.
	oldKey, oldCurrent, err := fetchCurrentRevisionOfSpec(ctx, client, request.GetName())
	if err == nil && oldCurrent != nil {
		oldCurrent.Currency = models.NotCurrent
		_, err = client.Put(ctx, oldKey, oldCurrent)
		if err != nil {
			log.Printf("oops %+v", err)
			return nil, internalError(err)
		}
	}
	// Make the selected revision the current revision by giving it a new RevisionID and saving it
	oldBlobKey := client.NewKey(models.BlobEntityName, spec.ResourceNameWithRevision())
	blob := &models.Blob{}
	err = client.Get(ctx, oldBlobKey, blob)
	if err != nil {
		return nil, internalError(err)
	}
	spec.BumpRevision()
	spec.Currency = models.IsCurrent
	newSpecKey := client.NewKey(models.SpecEntityName, spec.ResourceNameWithRevision())
	_, err = client.Put(ctx, newSpecKey, spec)
	if err != nil {
		return nil, internalError(err)
	}
	// Resave the blob for the current revision with the new RevisionID
	newBlobKey := client.NewKey(models.BlobEntityName, spec.ResourceNameWithRevision())
	blob.RevisionID = spec.RevisionID
	_, err = client.Put(ctx, newBlobKey, blob)
	if err != nil {
		return nil, internalError(err)
	}
	// Send a notification of the new revision.
	s.notify(rpc.Notification_UPDATED, spec.ResourceNameWithRevision())
	return spec.Message(nil, spec.RevisionID)
}

// fetchSpec gets the stored model of a Spec.
func fetchSpec(
	ctx context.Context,
	client storage.Client,
	name string,
) (*models.Spec, string, error) {
	spec, err := models.NewSpecFromResourceName(name)
	if err != nil {
		return nil, "", err
	}
	// if there's no revision, get the current revision
	if spec.RevisionID == "" {
		_, spec, err := fetchCurrentRevisionOfSpec(ctx, client, name)
		if err != nil {
			return nil, "", err
		}
		return spec, "", nil
	}
	// since a revision was specified, get the spec by revision
	// if the revision reference is a tag, resolve the tag
	var resourceName string
	var revisionName string
	specRevisionTag := &models.SpecRevisionTag{}
	k0 := client.NewKey(models.SpecRevisionTagEntityName, spec.ResourceNameWithRevision())
	err = client.Get(ctx, k0, specRevisionTag)
	if client.IsNotFound(err) {
		// if there is no tag, just use the provided revision
		resourceName = spec.ResourceNameWithRevision()
		revisionName = spec.RevisionID
	} else if err != nil {
		return nil, "", err
	} else {
		// if there is a tag, use the revision that the tag references
		resourceName = specRevisionTag.ResourceNameWithRevision()
		revisionName = specRevisionTag.Tag
	}
	// now that we know the revision, use it get the spec
	k := client.NewKey(models.SpecEntityName, resourceName)
	err = client.Get(ctx, k, spec)
	if client.IsNotFound(err) {
		return nil, revisionName, err
	} else if err != nil {
		return nil, revisionName, err
	}
	return spec, revisionName, nil
}

// fetchMostRecentNonCurrentRevisionOfSpec gets the most recent revision that's not current.
func (s *RegistryServer) fetchMostRecentNonCurrentRevisionOfSpec(
	ctx context.Context,
	client storage.Client,
	name string,
) (storage.Key, *models.Spec, error) {
	pattern, err := models.NewSpecFromResourceName(name)
	if err != nil {
		return nil, nil, err
	}
	// note that we ignore any specified RevisionID
	q := client.NewQuery(models.SpecEntityName)
	q = q.Require("ProjectID", pattern.ProjectID)
	q = q.Require("ApiID", pattern.ApiID)
	q = q.Require("VersionID", pattern.VersionID)
	q = q.Require("SpecID", pattern.SpecID)
	q = q.Require("Currency", models.NotCurrent)
	q = q.Order("-CreateTime")
	it := client.Run(ctx, q)

	if s.weTrustTheSort {
		spec := &models.Spec{}
		k, err := it.Next(spec)
		if err != nil {
			return nil, nil, client.NotFoundError()
		}
		return k, spec, nil
	} else {
		specs := make([]*models.Spec, 0)
		for {
			spec := &models.Spec{}
			_, err := it.Next(spec)
			if err != nil {
				break
			}
			specs = append(specs, spec)
		}
		sort.Slice(specs, func(i, j int) bool {
			return specs[i].CreateTime.After(specs[j].CreateTime)
		})
		k := client.NewKey("Spec", specs[0].Key)
		return k, specs[0], nil
	}
}

// fetchCurrentRevisionOfSpec gets the current revision.
func fetchCurrentRevisionOfSpec(
	ctx context.Context,
	client storage.Client,
	name string,
) (storage.Key, *models.Spec, error) {
	pattern, err := models.NewSpecFromResourceName(name)
	if err != nil {
		return nil, nil, err
	}
	// note that we ignore any specified RevisionID
	q := client.NewQuery(models.SpecEntityName)
	q = q.Require("ProjectID", pattern.ProjectID)
	q = q.Require("ApiID", pattern.ApiID)
	q = q.Require("VersionID", pattern.VersionID)
	q = q.Require("SpecID", pattern.SpecID)
	q = q.Require("Currency", models.IsCurrent)
	it := client.Run(ctx, q)
	spec := &models.Spec{}
	k, err := it.Next(spec)
	if err != nil {
		return nil, nil, client.NotFoundError()
	}
	return k, spec, nil
}

// fetchBlobForSpec gets the blob containing the spec contents.
func fetchBlobForSpec(
	ctx context.Context,
	client storage.Client,
	spec *models.Spec) (*models.Blob, error) {
	blob := &models.Blob{}
	k := client.NewKey(models.BlobEntityName, spec.ResourceNameWithRevision())
	err := client.Get(ctx, k, blob)
	return blob, err
}
