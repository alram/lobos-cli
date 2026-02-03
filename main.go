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
	nameFlag = &cli.StringFlag{
		Name:     "name",
		Usage:    "name of the new resource",
		Required: true,
	}
	forceFlag = &cli.BoolFlag{
		Name: "force",
	}
	accessKeyFlag = &cli.StringFlag{
		Name:     "access-key",
		Usage:    "S3 access key",
		Required: false,
	}
	accessKeyReqFlag = &cli.StringFlag{
		Name:     "access-key",
		Usage:    "S3 access key",
		Required: true,
	}
	secretFlag = &cli.StringFlag{
		Name:     "secret",
		Usage:    "S3 secret",
		Required: false,
	}

	rmUserFlags = []cli.Flag{
		nameFlag,
		forceFlag,
	}
	userFlags = []cli.Flag{
		nameFlag,
		accessKeyFlag,
		secretFlag,
	}
	keyRmFlags = []cli.Flag{
		nameFlag,
		accessKeyReqFlag,
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
						Flags: userFlags,
						Action: func(ctx context.Context, cmd *cli.Command) error {
							r, err := client.AddUser(ctx, &pb.User{
								Name:   cmd.String("name"),
								Key:    cmd.String("access-key"),
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
					{
						Name:  "rm",
						Usage: "Remove an user",
						Flags: rmUserFlags,
						Action: func(ctx context.Context, cmd *cli.Command) error {
							_, err := client.RmUser(ctx, &pb.RmUserParams{
								Name:  cmd.String("name"),
								Force: cmd.Bool("force"),
							})
							if err != nil {
								log.Fatalf("Could not list all users: %v", err)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "keys",
				Usage: "Add/Remove S3 access",
				Commands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Adds an S3 access",
						Flags: userFlags,
						Action: func(ctx context.Context, cmd *cli.Command) error {
							r, err := client.AddKey(ctx, &pb.User{
								Name:   cmd.String("name"),
								Key:    cmd.String("access-key"),
								Secret: cmd.String("secret"),
							})
							if err != nil {
								log.Fatalf("could not add user: %v", err)
							}

							fmt.Printf("Access: %s, secret: %s\n", r.User.Key, r.User.Secret)

							return nil
						},
					},
					{
						Name:  "rm",
						Usage: "Removes an S3 access",
						Flags: keyRmFlags,
						Action: func(ctx context.Context, cmd *cli.Command) error {
							_, err := client.RmKey(ctx, &pb.User{
								Name: cmd.String("name"),
								Key:  cmd.String("access-key"),
							})
							if err != nil {
								log.Fatalf("could not add user: %v", err)
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
