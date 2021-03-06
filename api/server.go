package api

// LICENSE: Apache
// This is taken from https://github.com/philips/grpc-gateway-example/
// PLEASE DO NOT call any config.* stuff here as it creates dependencies that can't be injected easily when testing

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/centrifuge/go-centrifuge/centerrors"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	logging "github.com/ipfs/go-log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logging.Logger("cent-api-server")

// CentAPIServer is an implementation of node.Server interface for serving HTTP based Centrifuge API
type CentAPIServer struct {
	Address     string
	Port        int
	CentNetwork string
}

func NewCentAPIServer(
	address string,
	port int,
	centNetwork string,
) *CentAPIServer {
	return &CentAPIServer{
		Address:     address,
		Port:        port,
		CentNetwork: centNetwork,
	}
}

func (*CentAPIServer) Name() string {
	return "CentAPIServer"
}

// Serve exposes the client APIs for interacting with a centrifuge node
func (c *CentAPIServer) Start(ctx context.Context, wg *sync.WaitGroup, startupErr chan<- error) {
	defer wg.Done()
	certPool, err := loadCertPool()
	if err != nil {
		startupErr <- err
	}
	keyPair, err := loadKeyPair()
	if err != nil {
		startupErr <- err
	}
	addr := c.Address

	creds := credentials.NewTLS(&tls.Config{
		RootCAs:            certPool,
		ServerName:         addr,
		Certificates:       []tls.Certificate{keyPair},
		InsecureSkipVerify: true,
	})

	opts := []grpc.ServerOption{grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName:         addr,
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()

	err = registerServices(ctx, grpcServer, gwmux, addr, dopts)
	if err != nil {
		startupErr <- err
		return
	}

	mux.Handle("/", gwmux)
	srv := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			NextProtos:   []string{"h2"},
		},
	}

	startUpErrOut := make(chan error)
	go func(startUpErrInner chan<- error) {
		conn, err := net.Listen("tcp", c.Address)
		if err != nil {
			startUpErrInner <- err
			return
		}
		log.Infof("HTTP/gRpc listening on Port: %d\n", c.Port)
		log.Infof("Connecting to Network: %s\n", c.CentNetwork)
		err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
		if err != nil {
			startUpErrInner <- err
		}
	}(startUpErrOut)

	// listen to context events as well as http server startup errors
	for {
		select {
		case err := <-startUpErrOut:
			// this could create an issue if the listeners are blocking.
			// We need to only propagate the error if its an error other than a server closed
			if err != nil && err.Error() != http.ErrServerClosed.Error() {
				startupErr <- err
				return
			}
			// most probably a graceful shutdown
			log.Info(err)
			return
		case <-ctx.Done():
			ctxn, _ := context.WithTimeout(context.Background(), 1*time.Second)
			// gracefully shutdown the server
			// we can only do this because srv is thread safe
			log.Info("Shutting down API server")
			err := srv.Shutdown(ctxn)
			if err != nil {
				panic(err)
			}
			log.Info("API server stopped")
			return
		}
	}

}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func loadCertPool() (certPool *x509.CertPool, err error) {
	certPool = x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(InsecureCert))
	if !ok {
		return nil, centerrors.Wrap(errors.New("could not load certpool"), "")
	}
	return certPool, nil
}

func loadKeyPair() (keyPair tls.Certificate, err error) {
	pair, err := tls.X509KeyPair([]byte(InsecureCert), []byte(InsecureKey))
	if err != nil {
		return pair, err
	}
	return pair, nil
}
