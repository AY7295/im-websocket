package service

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"webSocket-be/model"
	"webSocket-be/proto"
)

var userClient proto.UserClient
var ctx context.Context

func InitGRPC() error {

	ctx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"app-id":     viper.GetString("grpc.app_id"),
		"app-secret": viper.GetString("grpc.app_secret"),
	}))

	addr := flag.String(viper.GetString("grpc.name"), viper.GetString("grpc.address"), viper.GetString("grpc.usage"))
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
	user, err := userClient.VerifyToken(ctx, stringWrap)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Name: user.Name,
		Id:   user.Xh,
	}, nil

}
