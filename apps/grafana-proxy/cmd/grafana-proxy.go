// SPDX-FileCopyrightText: (C) 2025 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	proto "github.com/open-edge-platform/o11y-tenant-controller/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type proxy struct {
	mutex                 sync.RWMutex
	projects              []string
	includeDeleted        bool
	includeEdgenodeSystem bool

	mimirProxy *httputil.ReverseProxy
	lokiProxy  *httputil.ReverseProxy

	logger *zap.SugaredLogger
}

var globalRoles = []string{"admin", "viewer"}

const (
	errHTTPPermissionDenied = "Permission denied"
	errHTTPNoProjectsFound  = "No projects found"
)

func main() {
	var logLevel, mimirURL, lokiURL, otcURL string
	var includeDeleted, includeEdgenodeSystem bool
	var port int
	flag.StringVar(&logLevel, "log-level", "info", "set log level: debug, info, warn, error")
	flag.StringVar(&mimirURL, "mimir-url", "http://edgenode-observability-mimir-gateway.orch-infra.svc.cluster.local:8181/prometheus", "set mimir URL")
	flag.StringVar(&lokiURL, "loki-url", "http://edgenode-observability-loki-gateway.orch-infra.svc.cluster.local:80", "set loki URL")
	flag.StringVar(&otcURL, "otc-url", "observability-tenant-controller.orch-platform.svc.cluster.local:50051",
		"set observability tenant controller URL")
	flag.BoolVar(&includeDeleted, "include-deleted", false, "include deleted projects")
	flag.BoolVar(&includeEdgenodeSystem, "include-edgenode-system", false, "include edgenode-system project in Header")
	flag.IntVar(&port, "port", 9190, "set http server port")

	flag.Parse()

	logger, err := initLogger(logLevel)
	if err != nil {
		log.Panicf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Errorf("Failed to sync logger: %v", err)
		}
	}()

	mimir, err := url.Parse(mimirURL)
	if err != nil {
		logger.Panicf("Failed to parse mimir URL: %v", err)
	}

	loki, err := url.Parse(lokiURL)
	if err != nil {
		logger.Panicf("Failed to parse loki URL: %v", err)
	}

	proxy := &proxy{
		logger:                logger,
		mimirProxy:            httputil.NewSingleHostReverseProxy(mimir),
		lokiProxy:             httputil.NewSingleHostReverseProxy(loki),
		includeDeleted:        includeDeleted,
		includeEdgenodeSystem: includeEdgenodeSystem,
	}

	tenantControllerConn, err := grpc.NewClient(otcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.DefaultConfig}),
	)
	if err != nil {
		logger.Panicf("Failed to connect to tenant controller: %v", err)
	}
	defer tenantControllerConn.Close()

	http.HandleFunc("/", proxy.dispatcher())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go proxy.fetchProjectUpdates(ctx, tenantControllerConn)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Infof("Proxy server running on port: %v", port)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Panicf("Error while serving: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Errorf("Proxy server shutdown error: %v", err)
	}
}

func (p *proxy) proxyHandler(reverseProxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.logger.Debugf("Request received: %v", r.URL.Path)

		// Avoid logging the full HTTP headers. Instead, log only non-sensitive metadata, such as specific header names.
		headerNames := make([]string, 0, len(r.Header))
		for name := range r.Header {
			headerNames = append(headerNames, name)
		}
		p.logger.Debugf("Request header names: %v", headerNames)

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			p.logger.Errorf("Authorization header not found")
			http.Error(w, errHTTPPermissionDenied, http.StatusUnauthorized)
			return
		}
		// TODO: use short lived cache for token and roles to skip parsing token on every request
		realmRoles, resourceRoles, err := getTelemetryClientRolesFromJWT(token)
		if err != nil {
			p.logger.Errorf("Failed to parse token: %v", err)
			http.Error(w, errHTTPPermissionDenied, http.StatusUnauthorized)
			return
		}

		header, err := p.setHeaders(realmRoles, resourceRoles)
		if errors.Is(err, errNoProjectAccess) {
			p.logger.Error(err)
			http.Error(w, errHTTPPermissionDenied, http.StatusUnauthorized)
			return
		} else if err != nil {
			p.logger.Errorf("Failed to set headers: %v", err)
			http.Error(w, errHTTPNoProjectsFound, http.StatusInternalServerError)
			return
		}
		p.logger.Debugf("Header: %v", header)
		r.Header.Set("X-Scope-OrgID", header)

		reverseProxy.ServeHTTP(w, r)
	}
}

func (p *proxy) dispatcher() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/api/v1") {
			p.proxyHandler(p.mimirProxy)(w, r)
		} else if strings.HasPrefix(path, "/loki") {
			p.proxyHandler(p.lokiProxy)(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

var errNoProjectAccess = errors.New("user does not have access to any projects")

func (p *proxy) setHeaders(realmRoles, resourceRoles []string) (string, error) {
	projectsCopy := make([]string, len(p.projects))
	// make a copy of the projects list to shorten the lock time
	p.mutex.RLock()
	copy(projectsCopy, p.projects)
	p.mutex.RUnlock()

	if len(projectsCopy) == 0 {
		return "", errors.New("projects not found")
	}

	// Global roles that allow access to all projects are: "admin" and "viewer".
	// If the role is one of the global roles, return all projects.
	// If the role ends with a "tc-r" - it's a restricted viewer role to certain projects.
	isGlobalRole := false
	for _, role := range resourceRoles {
		if slices.Contains(globalRoles, role) {
			isGlobalRole = true
			break
		}
	}

	allowedProjects := make([]string, 0, len(projectsCopy)+1)
	if isGlobalRole {
		allowedProjects = append(allowedProjects, projectsCopy...)
		if p.includeEdgenodeSystem {
			allowedProjects = append(allowedProjects, "edgenode-system")
		}
	} else {
		allowedProjects = handleViewerRole(projectsCopy, realmRoles)
	}

	if len(allowedProjects) == 0 {
		return "", errNoProjectAccess
	}

	return strings.Join(allowedProjects, "|"), nil
}

func handleViewerRole(projects []string, roles []string) []string {
	projectSet := make(map[string]struct{}, len(projects))
	for _, project := range projects {
		projectSet[project] = struct{}{}
	}
	allowedProjects := make([]string, 0, len(roles))

	for _, role := range roles {
		if strings.HasSuffix(role, "_tc-r") {
			project := strings.TrimSuffix(role, "_tc-r")
			if _, exists := projectSet[project]; exists {
				allowedProjects = append(allowedProjects, project)
			}
		}
	}
	return allowedProjects
}

func (p *proxy) fetchProjectUpdates(ctx context.Context, tcConn *grpc.ClientConn) {
	tc := proto.NewProjectServiceClient(tcConn)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Context done, exiting project updates fetcher")
			return
		default:
			stream, err := tc.StreamProjectUpdates(ctx, &proto.EmptyRequest{})
			if err != nil {
				p.logger.Errorf("Failed to stream project updates: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for {
				update, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					p.logger.Info("Stream terminated gracefully, reconnecting")
					time.Sleep(100 * time.Millisecond)
					break
				} else if err != nil {
					p.logger.Errorf("Failed to receive update: %v", err)
					time.Sleep(1 * time.Second)
					break
				}
				prjs := make([]string, 0, len(update.GetProjects()))
				for _, project := range update.GetProjects() {
					if p.includeDeleted || project.Data.Status == "Created" {
						prjs = append(prjs, project.GetKey())
					}
				}

				p.mutex.Lock()
				p.projects = prjs
				p.mutex.Unlock()
				p.logger.Debugf("Update received, projects set to: %v", p.projects)
			}
		}
	}
}

func getTelemetryClientRolesFromJWT(token string) (realmRoles, resourceRoles []string, err error) {
	claims := jwt.MapClaims{}
	// ParseUnverified used as JWT Token is already verified (twice) on ingress level and on Grafana level to login.
	_, _, err = new(jwt.Parser).ParseUnverified(token, &claims)
	if err != nil {
		return nil, nil, err
	}

	realmRoles = collectRoles(claims, "realm_access", "roles")

	// Extract resource_access roles if present
	if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
		resourceRoles = collectRoles(resourceAccess, "telemetry-client", "roles")
	}

	return realmRoles, resourceRoles, nil
}

func collectRoles(data map[string]interface{}, mainKey, rolesKey string) []string {
	nested, ok := data[mainKey].(map[string]interface{})
	if !ok {
		return nil
	}
	rawRoles, ok := nested[rolesKey].([]interface{})
	if !ok {
		return nil
	}

	roles := make([]string, 0, len(rawRoles))
	for _, role := range rawRoles {
		if roleStr, ok := role.(string); ok {
			roles = append(roles, roleStr)
		}
	}
	return roles
}

func initLogger(logLevel string) (*zap.SugaredLogger, error) {
	log.Printf("Initializing logger with level: %v", logLevel)
	config := zap.Config{
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(level)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
