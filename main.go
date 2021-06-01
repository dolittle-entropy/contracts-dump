package main

import (
	"flag"
	"net"

	grpcproxy "github.com/bradleyjkemp/grpc-tools/grpc-proxy"
	"github.com/sirupsen/logrus"
)

var (
	listen  string
	runtime string
)

func main() {
	flag.StringVar(&listen, "listen", "0.0.0.0:50053", "The address to listen to")
	flag.StringVar(&runtime, "runtime", "localhost:50553", "The address of the Runtime to proxy to")
	flag.Parse()

	grpcproxy.RegisterDefaultFlags()

	listenHost, listenPort, err := net.SplitHostPort(listen)
	if err != nil {
		logrus.WithError(err).Fatal("Could not parse address to listen to")
	}

	flag.Set("interface", listenHost)
	flag.Set("port", listenPort)
	flag.Set("destination", runtime)

	parser, err := NewParser("contracts/Source")
	if err != nil {
		logrus.WithError(err).Fatal("Could not load Contracts .proto files")
	}

	interceptor := ContractsInterceptor(parser)

	proxy, err := grpcproxy.New(
		grpcproxy.DefaultFlags(),
		grpcproxy.WithInterceptor(interceptor),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create proxy")
		return
	}
	logrus.Infof("Forwarding to Runtime on %s", runtime)

	if err := proxy.Start(); err != nil {
		logrus.WithError(err).Fatal("Failed to start proxy")
	}
}
