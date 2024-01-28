module github.com/kkiling/photo-library/backend/photo_sync

go 1.20

require (
	github.com/cheggaaa/pb/v3 v3.1.4
	github.com/hirochachacha/go-smb2 v1.1.0
	github.com/kkiling/photo-library/backend/api v0.0.0
	github.com/mattn/go-sqlite3 v1.14.17
	golang.org/x/crypto v0.12.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/geoffgarside/ber v1.1.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230726155614-23370e0ffb3e // indirect
)

replace github.com/kkiling/photo-library/backend/api v0.0.0 => ../api
