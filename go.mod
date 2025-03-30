module github.com/garden-raccoon/delivery-pkg

go 1.23.2

require (
	github.com/goccy/go-json v0.10.5
	github.com/gofrs/uuid v4.4.0+incompatible
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

require (
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241104194629-dd2ea8efbc28 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.67.1
