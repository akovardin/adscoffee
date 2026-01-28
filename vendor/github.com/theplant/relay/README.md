# relay

`relay` is a library designed to simplify Relay-style pagination in Go applications, supporting both keyset-based and offset-based pagination. It helps developers efficiently implement pagination queries while offering optimization options, such as skipping `TotalCount` queries and encrypting cursors.

**Currently, only the `GORM` adapter is provided by default. For other adapters, refer to `gormrelay` to implement them yourself.**

## Features

- **Supports keyset-based and offset-based pagination**: You can freely choose high-performance keyset pagination based on multiple indexed columns, or use offset pagination.
- **Optional cursor encryption**: Supports encrypting cursors using `GCM(AES)` or `Base64` to ensure the security of pagination information.
- **Flexible query strategies**: Optionally skip the `TotalCount` query to improve performance, especially in large datasets.
- **Non-generic support**: Even without using Go generics, you can paginate using the `any` type for flexible use cases.

## Usage

### Basic Usage

```go
p := relay.New(
    cursor.Base64(func(ctx context.Context, req *relay.ApplyCursorsRequest) (*relay.ApplyCursorsResponse[*User], error) {
        // Offset-based pagination
        // return gormrelay.NewOffsetAdapter[*User](db)(ctx, req)

        // Keyset-based pagination
        return gormrelay.NewKeysetAdapter[*User](db)(ctx, req)
    }),
    // defaultLimit / maxLimit
    relay.EnsureLimits[*User](10, 100),
    // Append primary sorting fields, if any are unspecified
    relay.EnsurePrimaryOrderBy[*User](
        relay.OrderBy{Field: "ID", Desc: false},
        relay.OrderBy{Field: "Version", Desc: false},
    ),
)

conn, err := p.Paginate(
    context.Background(),
    // relay.WithSkip(context.Background(), relay.Skip{
    //     Edges:      true,
    //     Nodes:      true,
    //     PageInfo:   true,
    //     TotalCount: true,
    // }),

    // Query first 10 records
    &relay.PaginateRequest[*User]{
        First: lo.ToPtr(10),
    }
)
```

### Cursor Encryption

If you need to encrypt cursors, you can use `cursor.Base64` or `cursor.GCM` wrappers:

```go
// Encode cursors with Base64
cursor.Base64(gormrelay.NewOffsetAdapter[*User](db))

// Encrypt cursors with GCM(AES)
gcm, err := cursor.NewGCM(encryptionKey)
require.NoError(t, err)
cursor.GCM(gcm)(gormrelay.NewKeysetAdapter[*User](db))
```

### Non-Generic Usage

If you do not use generics, you can create a paginator with the `any` type and combine it with the `db.Model` method:

```go
p := relay.New(
    func(ctx context.Context, req *relay.ApplyCursorsRequest) (*relay.ApplyCursorsResponse[any], error) {
        // Since this is a generic function (T: any), we must call db.Model(x)
        return gormrelay.NewKeysetAdapter[any](db.Model(&User{}))(ctx, req)
    },
    relay.EnsureLimits[any](10, 100),
    relay.EnsurePrimaryOrderBy[any](relay.OrderBy{Field: "ID", Desc: false}),
)
conn, err := p.Paginate(context.Background(), &relay.PaginateRequest[any]{
    First: lo.ToPtr(10), // query first 10 records
})
```

## Reference

For more information about Relay-style pagination, refer to [GraphQL Connections](https://relay.dev/graphql/connections.htm).
