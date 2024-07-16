package main

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    token := os.Getenv("AUTHZED_TOKEN")
    if token == "" {
        log.Fatalf("AUTHZED_TOKEN not set in .env file")
    }

    client, err := newSpiceDBClient(token)
    if err != nil {
        log.Fatalf("Failed to create SpiceDB client: %v", err)
    }

    testCases := []struct {
        user       string
        bucket     string
        permission string
        expected   bool
    }{
        {"alice", "test-bucket", "read", true},
        {"alice", "test-bucket", "write", true},
        {"alice", "test-bucket", "admin", true},
        {"bob", "test-bucket", "read", true},
        {"bob", "test-bucket", "write", false},
        {"bob", "test-bucket", "admin", false},
        {"charlie", "test-bucket", "read", true},
    }

    for _, tc := range testCases {
        hasPermission, err := checkPermission(client, tc.user, tc.bucket, tc.permission)
        if err != nil {
            log.Printf("Error checking permission for %s on %s: %v", tc.user, tc.bucket, err)
            continue
        }

        if hasPermission == tc.expected {
            fmt.Printf("✅ User %s %s permission %s on %s\n", tc.user, map[bool]string{true: "has", false: "doesn't have"}[hasPermission], tc.permission, tc.bucket)
        } else {
            fmt.Printf("❌ Unexpected result: User %s %s permission %s on %s\n", tc.user, map[bool]string{true: "has", false: "doesn't have"}[hasPermission], tc.permission, tc.bucket)
        }
    }
}
