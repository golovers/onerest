# Simple Version One REST Client

This is a very simple REST client for Version One

## Examples

```go
builder := onerest.OneRestBuilder{}
// With access token
service := builder.Host("host url").WithAccessToken("token")
// Or with username password
service := builder.Host(config.Host).WithUserPassword("username", "password")

// Query data from remote server
service.Scope("exampleScopeName").Stories(...)
service.Scope("exampleScopeName").Defects(...)
....
```

