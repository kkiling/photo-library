package main

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter/smbread"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter/sqlitestorage"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/syncfiles"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	smb := smbread.NewSmbRead(smbread.Config{
		User:     "guest",
		Password: "guest",
		Address:  "10.10.10.204:445",
		// Address:    "nas.lan:445",
		ShareName:  "storage",
		DirPath:    "photos",
		Extensions: []string{".jpg", ".jpeg", ".png", ".bmp"},
	})

	err := smb.Connect(ctx)
	if err != nil {
		fmt.Printf("fail smb.Connect: %v", err)
		return
	}
	defer func(smb *smbread.SmbRead) {
		if err := smb.Disconnect(); err != nil {
			fmt.Printf("fail smb.Disconnect: %v", err)
		}
	}(smb)

	storage, err := sqlitestorage.NewStorage(sqlitestorage.Config{
		DSN: "./files.db",
	})

	if err != nil {
		fmt.Printf("fail sqlitestorage.NewStorage: %v", err)
		return
	}

	var sync = syncfiles.NewSyncPhotos(smb, storage, syncfiles.Config{
		GrpcServerHost: "localhost:8181",
		ClientId:       "mbp-kkiling-OZON-HXW066MJFG",
		NumWorkers:     8,
	})
	if err := sync.Sync(ctx); err != nil {
		fmt.Printf("fail sync.Sync: %v", err)
		return
	}
}
