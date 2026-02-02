package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "lobos/lobosctl/protos"

	"github.com/urfave/cli/v3"
)

var (
	globalFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "endpoint",
			Usage: "gRPC server endpoint",
			Value: "127.0.0.1:50051",
		},
	}
	nameFlag = []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "name of the new user",
			Required: true,
		},
	}
	fullUsersFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "name of the new user",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "access",
			Usage:    "S3 access key",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "secret",
			Usage:    "S3 secret",
			Required: false,
		},
	}
)

var (
	conn   *grpc.ClientConn
	client pb.ControlPlaneClient
)

func main() {

	cmd := &cli.Command{
		Name:  "lobosctl",
		Usage: "Interact with lobos control plane",
		Flags: globalFlags,
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			endpoint := cmd.String("endpoint")

			var err error
			conn, err = grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return ctx, fmt.Errorf("did not connect: %v", err)
			}

			client = pb.NewControlPlaneClient(conn)

			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:  "users",
				Usage: "Add/Remove/List S3 users",
				Commands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Adds an user",
						Flags: fullUsersFlags,
						Action: func(ctx context.Context, cmd *cli.Command) error {
							r, err := client.AddUser(ctx, &pb.User{
								Name:   cmd.String("name"),
								Key:    cmd.String("access"),
								Secret: cmd.String("secret"),
							})
							if err != nil {
								log.Fatalf("could not add user: %v", err)
							}

							fmt.Printf("User %s created. Access: %s, secret: %s\n", r.User.Name, r.User.Key, r.User.Secret)

							return nil
						},
					},
					{
						Name:  "list",
						Usage: "List all users",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							r, err := client.ListAllUsers(ctx, &pb.Filters{})
							if err != nil {
								log.Fatalf("Could not list all users: %v", err)
							}
							for _, u := range r.Users {
								fmt.Printf("user: %s - key: %s - secret: %s\n", u.User.Name, u.User.Key, u.User.Secret)
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
