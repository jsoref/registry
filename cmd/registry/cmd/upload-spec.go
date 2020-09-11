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

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/gapic"
	rpcpb "github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
)

func init() {
	uploadCmd.AddCommand(specCmd)
	specCmd.Flags().String("version", "", "Resource name of version for uploaded spec")
}

// specCmd represents the spec command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Upload files of an API spec.",
	Long:  "Upload files of an API spec.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()
		flagset := cmd.LocalFlags()
		version, err := flagset.GetString("version")
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		fmt.Printf("spec called with args %+v and version %s\n", args, version)

		client, err := connection.NewClient(ctx)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		for _, arg := range args {
			matches, err := filepath.Glob(arg)
			if err != nil {
				log.Printf("%s\n", err.Error())
			}
			// for each match, upload the file
			for _, match := range matches {
				log.Printf("now upload %+v", match)
				fi, err := os.Stat(match)
				if err == nil {
					switch mode := fi.Mode(); {
					case mode.IsDir():
						fmt.Printf("upload directory %s\n", match)
						uploadDirectory(match, client, version)
					case mode.IsRegular():
						fmt.Printf("upload file %s\n", match)
						uploadSpecFile(match, client, version)
					}
				} else {
					log.Printf("%+v", err)
				}
			}
		}
	},
}

func uploadDirectory(dirname string, client *gapic.RegistryClient, version string) error {
	return filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			log.Printf("%+s", path)
			if err != nil {
				return err
			}
			if !info.IsDir() {
				uploadSpecFile(path, client, version)
			}
			return nil
		})
}

func uploadSpecFile(filename string, client *gapic.RegistryClient, version string) {
	// does the spec file exist? if not, create it
	request := &rpcpb.GetSpecRequest{}
	request.Name = version + "/specs/" + filename
	ctx := context.TODO()
	response, err := client.GetSpec(ctx, request)
	log.Printf("response %+v\nerr %+v", response, err)
	if err != nil { // TODO only do this for NotFound errors
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("err %+v", err)
		} else {
			request := &rpcpb.CreateSpecRequest{}
			request.Parent = version
			request.SpecId = filename
			request.Spec = &rpcpb.Spec{}
			request.Spec.Filename = filename
			request.Spec.Contents = bytes
			switch filename {
			case "swagger.yaml":
				request.Spec.Style = "openapi-v2"
			case "openapi.yaml":
				request.Spec.Style = "openapi-v3"
			default:
				request.Spec.Style = "proto"
			}
			response, err := client.CreateSpec(ctx, request)
			log.Printf("response %+v\nerr %+v", response, err)
		}
	}
}
