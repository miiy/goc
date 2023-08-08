package server

//
//import (
//	"context"
//	pb "docxlib.com/api/auth/v1"
//	"docxlib.com/pkg/log"
//	"docxlib.com/pkg/registry"
//	"docxlib.com/pkg/server/grpc"
//	"docxlib.com/service/auth/internal/server"
//	"flag"
//)
//
//func main() {
//	conf := flag.String("c", "./configs/default.yaml", "config file")
//	etcd := flag.String("etcd", "127.0.0.1:2379", "the etcd endpoints")
//	flag.Parse()
//
//	app, cleanUp, err := newApp(*conf)
//	if err != nil {
//		panic(err)
//	}
//	defer cleanUp()
//
//	s, err := grpc.NewServer(app.Config.Server.Grpc.Addr, app.Logger.Logger)
//	if err != nil {
//		app.Logger.Fatal("Failed to create server", log.Error(err))
//	}
//
//	pb.RegisterAuthServiceServer(s, server.NewAuthServiceServer(app.Database, app.JwtAuth))
//
//	registryConfig := registry.Config{
//		ServiceName:  app.Config.App.Name,
//		EtcdServer:   []string{*etcd},
//		InstanceAddr: app.Config.Server.Grpc.Addr,
//	}
//	r, err := registry.NewRegistry(&registryConfig, context.Background())
//	if err != nil {
//		panic(err)
//	}
//	cleanUp1 := r.Register()
//	defer cleanUp1()
//
//	if err = s.Serve(); err != nil {
//		app.Logger.Fatal("failed to serve: %v", log.Error(err))
//	}
//}
