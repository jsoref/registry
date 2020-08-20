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

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/gogo/googleapis/google/rpc"
)

var (
	grpcport = flag.String("grpcport", ":50051", "grpcport")
	conn     *grpc.ClientConn
	hs       *health.Server
)

const (
	address   string = ":50051"
	trustJWTs bool   = true // if true, JWT claims are accepted without signature verification.
)

// healthServer implements the gRPC health check service.
type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

// authorizationServer implements the Envoy authz service.
type authorizationServer struct{}

// Check implements the check operation in the Envoy authz service.
func (a *authorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Authorization called check()")

	b, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  ")
	if err == nil {
		log.Println("Inbound Headers: " + string(b))
	}
	ct, err := json.MarshalIndent(req.Attributes.ContextExtensions, "", "  ")
	if err == nil {
		log.Println("Context Extensions: " + string(ct))
	}

	authHeader, ok := req.Attributes.Request.Http.Headers["authorization"]
	if !ok {
		// there's no auth header, so the request is uncredentialed.
		return denyUncredentialedRequest(), nil
	}
	re := regexp.MustCompile("^[bB]earer[ ]+(.*)$")
	m := re.FindStringSubmatch(authHeader)
	if m == nil {
		return denyMalformedCredentials(), nil
	}
	credential := m[1]

	if isJWTToken(credential) {
		if trustJWTs {
			// get the user email from the token
			email := getJWTTokenEmail(credential)
			if email != "" {
				return allowAuthorizedUser(email), nil
			}
		}
		// try to verify an identity token
		token, err := getVerifiedToken(credential)
		if err == nil && token != nil {
			if isWriter(token.Email) || isReadOnlyMethod(req.Attributes.Request.Http.Headers[":path"]) {
				// the user is authorized so we allow the call.
				return allowAuthorizedUser(token.Email), nil
			}
			// we have a user, but they aren't authorized to do this.
			return denyUnauthorizedUser(), nil
		}
		if err != nil {
			log.Printf("%s", err.Error())
		}
	} else {
		// try to verify an access token
		user, err := getUser(credential)
		if err == nil && user != nil {
			if isWriter(user.Email) || isReadOnlyMethod(req.Attributes.Request.Http.Headers[":path"]) {
				// the user is authorized so we allow the call.
				return allowAuthorizedUser(user.Email), nil
			}
			// we have a user, but they aren't authorized to do this.
			return denyUnauthorizedUser(), nil
		}
		if err != nil {
			log.Printf("%s", err.Error())
		}
	}

	// we can't find a user for the auth header, so the user is unauthenticated.
	return denyUnauthenticatedUser(), nil
}

// isReadOnlyMethod recognizes Get and List operations as immutable.
func isReadOnlyMethod(path string) bool {
	methodName := filepath.Base(path)
	if strings.HasPrefix(methodName, "Get") ||
		strings.HasPrefix(methodName, "List") {
		return true
	}
	return false
}

// TODO: read this from a yaml file
var writers = []string{"timburks@google.com", "timburks@gmail.com"}

// isWriter returns true if a user is allowed to make mutable operations.
func isWriter(email string) bool {
	for _, writer := range writers {
		if email == writer {
			return true
		}
	}
	return false
}

type jwtTokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func isJWTToken(credential string) bool {
	parts := strings.Split(credential, ".")
	if len(parts) != 3 {
		return false
	}
	header := parts[0]
	v, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return false
	}
	var tokenHeader jwtTokenHeader
	json.Unmarshal(v, &tokenHeader)
	return tokenHeader.Typ == "JWT"
}

type jwtTokenPayload struct {
	Email string `json:"email"`
}

func getJWTTokenEmail(credential string) string {
	parts := strings.Split(credential, ".")
	if len(parts) != 3 {
		return ""
	}
	payload := parts[1]
	v, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return ""
	}
	log.Printf("token payload %+v\n", string(v))
	var tokenPayload jwtTokenPayload
	json.Unmarshal(v, &tokenPayload)
	return tokenPayload.Email
}

// GoogleUser holds information about a Google user.
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	PictureURL    string `json:"picture"`
}

// in-memory cache of users
var users map[string]*GoogleUser

func getUser(credential string) (*GoogleUser, error) {
	if users == nil {
		users = make(map[string]*GoogleUser)
	}
	// first check the cache
	cachedUser := users[credential]
	if cachedUser != nil {
		log.Printf("cached user: %+v for %s", cachedUser, credential)
		return cachedUser, nil
	}
	// otherwise, call the Google userinfo API
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+credential)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Unsuccessful response from userinfo service: %d (%s): %s",
			resp.StatusCode, resp.Status, string(b))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	user := &GoogleUser{}
	err = json.Unmarshal(b, user)
	if err != nil {
		return nil, err
	}
	users[credential] = user
	log.Printf("verified user: %+v for %s", user, credential)
	return user, nil
}

// GoogleToken holds information about a Google identity token (a JWT).
type GoogleToken struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
}

// in-memory cache of tokens
var tokens map[string]*GoogleToken

func getVerifiedToken(credential string) (*GoogleToken, error) {
	if tokens == nil {
		tokens = make(map[string]*GoogleToken)
	}
	// first check the cache
	cachedToken := tokens[credential]
	if cachedToken != nil {
		log.Printf("cached token: %+v for %s", cachedToken, credential)
		return cachedToken, nil
	}
	// otherwise, call the Google tokeninfo API
	req, err := http.NewRequest("GET", "https://oauth2.googleapis.com/tokeninfo", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("id_token", credential)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Unsuccessful response from tokeninfo service: %d (%s): %s",
			resp.StatusCode, resp.Status, string(b))
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	token := &GoogleToken{}
	err = json.Unmarshal(b, token)
	if err != nil {
		return nil, err
	}
	tokens[credential] = token
	log.Printf("verified token: %+v for %s", token, credential)
	return token, nil
}

func allowAuthorizedUser(username string) *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: []*core.HeaderValueOption{
					{
						Header: &core.HeaderValue{
							Key:   "x-authz-user",
							Value: username,
						},
					},
				},
			},
		},
	}
}

func denyUnauthorizedUser() *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.PERMISSION_DENIED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: "Permission denied",
			},
		},
	}
}

func denyUnauthenticatedUser() *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: "Authorization cannot be validated",
			},
		},
	}
}

func denyUncredentialedRequest() *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: "Authorization is missing",
			},
		},
	}
}

func denyMalformedCredentials() *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: "Authorization header is malformed",
			},
		},
	}
}

func main() {
	flag.Parse()

	if *grpcport == "" {
		fmt.Fprintln(os.Stderr, "missing -grpcport flag (:50051)")
		flag.Usage()
		os.Exit(2)
	}

	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
	s := grpc.NewServer(opts...)
	auth.RegisterAuthorizationServer(s, &authorizationServer{})
	healthpb.RegisterHealthServer(s, &healthServer{})

	log.Printf("Starting gRPC Server at %s", *grpcport)
	s.Serve(lis)
}
