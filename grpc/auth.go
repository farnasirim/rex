package grpc

import (
	"context"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/farnasirim/rex"
)

// AuthInfoInterceptor extracts peer's id from the certificate and adds it
// to the request context
func AuthInfoInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	grpcPeer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated,
			rex.ErrTLSCredentials.Error())
	}
	tlsInfo, ok := grpcPeer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated,
			rex.ErrTLSCredentials.Error())
	}

	if len(tlsInfo.State.PeerCertificates) == 0 {
		return nil, status.Errorf(codes.Unauthenticated,
			rex.ErrTLSCredentials.Error())
	} else if len(tlsInfo.State.PeerCertificates) > 1 {
		log.Warnln("Peer used multiple certificates. Using the first one.")
	}

	log.Debugln("CN:", tlsInfo.State.PeerCertificates[0].Subject.CommonName)
	ctx = context.WithValue(ctx,
		rex.UserIDContextKey,
		tlsInfo.State.PeerCertificates[0].Subject.CommonName,
	)

	log.Debugln(info.FullMethod)

	return handler(ctx, req)
}
