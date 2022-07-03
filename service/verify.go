package service

import (
	"context"
	"crypto/tls"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"webSocket-be/config"
	"webSocket-be/model"
	"webSocket-be/proto"
)

var userClient proto.UserClient
var ctx context.Context

func InitGRPC(conf config.Config) error {

	ctx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"app-id":     conf.AppID,
		"app-secret": conf.AppSecret,
	}))

	addr := flag.String(conf.RPCName, conf.RPCAddress, conf.RPCUsage)
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))
	if err != nil {
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(conn)
	userClient = proto.NewUserClient(conn)

	return nil
}

func verifyToken(token string) (*model.User, error) {

	stringWrap := &proto.StringWrap{Val: token}
	user, err1 := userClient.VerifyToken(ctx, stringWrap)
	if err1 != nil {
		return nil, err1
	}

	u := &model.User{
		Name: user.Name,
		Id:   user.Xh,
	}

	return u, nil
}
