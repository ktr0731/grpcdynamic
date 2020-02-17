package grpcdynamic

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type unaryMethod struct {
	name     string
	handler  grpc.UnaryHandler
	req, res interface{}
}

type Service struct {
	name    string
	methods []*unaryMethod
}

func NewService(name string) *Service {
	return &Service{
		name: name,
	}
}

func (s *Service) RegisterUnaryMethod(name string, req, res interface{}, h grpc.UnaryHandler) {
	s.methods = append(s.methods, &unaryMethod{name: name, handler: h, req: req, res: res})
}

func (s *Service) FullMethodName(methodName string) string {
	return fullMethod(s.name, methodName)
}

func NewServer(services []*Service, opt ...grpc.ServerOption) *grpc.Server {
	s := grpc.NewServer(opt...)
	for _, service := range services {
		serviceDesc := createServiceDesc(service)
		s.RegisterService(serviceDesc, struct{}{})
	}
	return s
}

func createServiceDesc(s *Service) *grpc.ServiceDesc {
	type emptyInterface interface{}

	sd := &grpc.ServiceDesc{
		ServiceName: s.name,
		HandlerType: (*emptyInterface)(nil),
		Methods:     make([]grpc.MethodDesc, 0, len(s.methods)),
		Metadata:    "api.proto",
	}
	for _, method := range s.methods {
		sd.Methods = append(sd.Methods, createMethodDesc(fullMethod(s.name, method.name), method))
	}
	return sd
}

func createMethodDesc(fullMethod string, m *unaryMethod) grpc.MethodDesc {
	return grpc.MethodDesc{
		MethodName: m.name,
		Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
			if err := dec(m.req); err != nil {
				return nil, err
			}
			if interceptor == nil {
				return m.handler(ctx, m.req)
			}
			info := &grpc.UnaryServerInfo{
				Server:     srv,
				FullMethod: fullMethod,
			}
			return interceptor(ctx, m.req, info, m.handler)
		},
	}
}

func fullMethod(serviceName, methodName string) string {
	return fmt.Sprintf("/%s/%s", serviceName, methodName)
}
