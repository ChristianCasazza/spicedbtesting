package main

import (
    "context"
    "fmt"
    "log"

    "github.com/authzed/authzed-go/v1"
    "github.com/authzed/grpcutil"
    pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
)

func newSpiceDBClient(token string) (*authzed.Client, error) {
    systemCerts, err := grpcutil.WithSystemCerts(grpcutil.VerifyCA)
    if err != nil {
        log.Fatalf("unable to load system CA certificates: %s", err)
    }

    client, err := authzed.NewClient(
        "grpc.authzed.com:443",
        systemCerts,
        grpcutil.WithBearerToken(token),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create SpiceDB client: %w", err)
    }
    return client, nil
}

func checkPermission(client *authzed.Client, userId, bucketId, permission string) (bool, error) {
    ctx := context.Background()

    req := &pb.CheckPermissionRequest{
        Resource: &pb.ObjectReference{
            ObjectType: "oceanprotocol_testing/bucket",
            ObjectId:   bucketId,
        },
        Permission: permission,
        Subject: &pb.SubjectReference{
            Object: &pb.ObjectReference{
                ObjectType: "oceanprotocol_testing/user",
                ObjectId:   userId,
            },
        },
    }

    resp, err := client.CheckPermission(ctx, req)
    if err != nil {
        return false, fmt.Errorf("failed to check permission: %w", err)
    }

    return resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}
