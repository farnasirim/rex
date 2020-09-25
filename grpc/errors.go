package grpc

import (
	"context"
	"encoding/json"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	//"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
)

type errorChain struct {
	Message string
	Next    *errorChain `json:",omitempty"`
}

func (e *errorChain) Is(other error) bool {
	return e.Error() == other.Error()
}

func (e *errorChain) Error() string {
	return e.Message
}

func (e *errorChain) Unwrap() error {
	return e.Next
}

func (e *errorChain) Marshal() error {
	marshalledError, err := json.Marshal(e)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return errors.New(string(marshalledError))
}

func errorChainFromError(err error) *errorChain {
	if err == nil {
		return nil
	}
	return &errorChain{
		Message: err.Error(),
		Next:    errorChainFromError(errors.Unwrap(err)),
	}
}

func errorChainFromJSON(marshalledError string) (*errorChain, error) {
	var chain errorChain
	if err := json.Unmarshal([]byte(marshalledError), &chain); err != nil {
		return nil, err
	}
	return &chain, nil
}

// ErrorUnarshallerInterceptor unmarshals the returned error from the grpc
// call into an error chain if it is the marshalled form of an error chain
func ErrorUnmarshallerInterceptor(
	ctx context.Context,
	method string,
	req,
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	topLevelErr := invoker(ctx, method, req, reply, cc, opts...)

	if topLevelErr != nil {
		st, _ := status.FromError(topLevelErr)
		if st == nil || st.Proto() == nil {
			return topLevelErr
		}
		errChain, err := errorChainFromJSON(st.Proto().GetMessage())
		if err != nil {
			return topLevelErr
		}
		return errChain
	}
	return topLevelErr
}

// ErrorMarshallerInterceptor marshals an error chain to json.
func ErrorMarshallerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	ret, err := handler(ctx, req)
	if st, ok := status.FromError(err); !ok {
		return ret, status.Errorf(st.Code(), errorChainFromError(err).Marshal().Error())
	}

	return ret, err
}
