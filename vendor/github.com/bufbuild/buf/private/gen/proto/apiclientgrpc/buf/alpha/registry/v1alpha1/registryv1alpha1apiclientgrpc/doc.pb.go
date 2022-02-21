// Copyright 2020-2022 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-apiclientgrpc. DO NOT EDIT.

package registryv1alpha1apiclientgrpc

import (
	context "context"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	zap "go.uber.org/zap"
)

type docService struct {
	logger          *zap.Logger
	client          v1alpha1.DocServiceClient
	contextModifier func(context.Context) context.Context
}

// GetSourceDirectoryInfo retrieves the directory and file structure for the
// given owner, repository and reference.
//
// The purpose of this is to get a representation of the file tree for a given
// module to enable exploring the module by navigating through its contents.
func (s *docService) GetSourceDirectoryInfo(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
) (root *v1alpha1.FileInfo, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.GetSourceDirectoryInfo(
		ctx,
		&v1alpha1.GetSourceDirectoryInfoRequest{
			Owner:      owner,
			Repository: repository,
			Reference:  reference,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.Root, nil
}

// GetSourceFile retrieves the source contents for the given owner, repository,
// reference, and path.
func (s *docService) GetSourceFile(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
	path string,
) (content []byte, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.GetSourceFile(
		ctx,
		&v1alpha1.GetSourceFileRequest{
			Owner:      owner,
			Repository: repository,
			Reference:  reference,
			Path:       path,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.Content, nil
}

// GetModulePackages retrieves the list of packages for the module based on the given
// owner, repository, and reference.
func (s *docService) GetModulePackages(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
) (name string, modulePackages []*v1alpha1.ModulePackage, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.GetModulePackages(
		ctx,
		&v1alpha1.GetModulePackagesRequest{
			Owner:      owner,
			Repository: repository,
			Reference:  reference,
		},
	)
	if err != nil {
		return "", nil, err
	}
	return response.Name, response.ModulePackages, nil
}

// GetModuleDocumentation retrieves the documentation for module based on the given
// owner, repository, and reference.
func (s *docService) GetModuleDocumentation(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
) (moduleDocumentation *v1alpha1.ModuleDocumentation, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.GetModuleDocumentation(
		ctx,
		&v1alpha1.GetModuleDocumentationRequest{
			Owner:      owner,
			Repository: repository,
			Reference:  reference,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.ModuleDocumentation, nil
}

// GetPackageDocumentation retrieves a a slice of documentation structures
// for the given owner, repository, reference, and package name.
func (s *docService) GetPackageDocumentation(
	ctx context.Context,
	owner string,
	repository string,
	reference string,
	packageName string,
) (packageDocumentation *v1alpha1.PackageDocumentation, _ error) {
	if s.contextModifier != nil {
		ctx = s.contextModifier(ctx)
	}
	response, err := s.client.GetPackageDocumentation(
		ctx,
		&v1alpha1.GetPackageDocumentationRequest{
			Owner:       owner,
			Repository:  repository,
			Reference:   reference,
			PackageName: packageName,
		},
	)
	if err != nil {
		return nil, err
	}
	return response.PackageDocumentation, nil
}
