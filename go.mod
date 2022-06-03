module github.com/aaarafat/vanessa

go 1.18

require (
	github.com/cornelk/hashmap v1.0.1
	github.com/mdlayher/packet v1.0.0
)

require (
	github.com/dchest/siphash v1.1.0 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/josharian/native v1.0.0 // indirect
	github.com/mdlayher/socket v0.2.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/text v0.3.3 // indirect
)

require (
	github.com/AkihiroSuda/go-netfilter-queue v0.0.0-20210408043041-d1559d12dfd4
	github.com/mdlayher/ethernet v0.0.0-20220221185849-529eae5b6118
	google.golang.org/grpc v1.47.0
	google.golang.org/grpc/examples v0.0.0-20220602231701-13b378bc4585
	google.golang.org/protobuf v1.28.0 // indirect
)

require (
	github.com/golang/protobuf v1.5.2
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	protos v0.0.0-00010101000000-000000000000
)

replace protos => ./apps/scripts/gRPC/protos
