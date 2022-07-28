package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"webSocket-be/config"
	"webSocket-be/model"
	"webSocket-be/proto"
)

func VerifyToken(token string) (*model.User, error) {

	conn, err := grpc.Dial(viper.GetString("grpc.address"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))
	if err != nil {
		return nil, fmt.Errorf("connect to grpc err: %w", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			config.Logfile.Println(fmt.Errorf("connect to grpc err: %w", err))
			return
		}
	}(conn)
	userClient := proto.NewUserClient(conn)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"app-id":     viper.GetString("grpc.app_id"),
		"app-secret": viper.GetString("grpc.app_secret"),
	}))

	stringWrap := &proto.StringWrap{Val: token}
	user, err := userClient.VerifyToken(ctx, stringWrap)
	if err != nil {
		config.Logfile.Println(fmt.Errorf("grpc verify token err: %w", err))
		return GetUserByToken(token)
	}

	return &model.User{
		Name: user.Name,
		Id:   user.Xh,
	}, nil

}
