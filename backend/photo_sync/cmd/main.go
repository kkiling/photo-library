package main

import (
	"context"
	"fmt"

	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter/smbread"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter/sqlitestorage"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/adapter/syncclient"
	"github.com/kkiling/photo-library/backend/photo_sync/internal/syncfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SmbConfig struct {
	User       string   `yaml:"user"`
	Password   string   `yaml:"password"`
	Address    string   `yaml:"address"`
	ShareName  string   `yaml:"shareName"`
	DirPath    string   `yaml:"dirPath"`
	Extensions []string `yaml:"extensions"`
}

type SyncConfig struct {
	SmbConfig              *SmbConfig `yaml:"smbConfig"`
	SqliteFile             string     `yaml:"sqliteFile"`
	ClientID               string     `yaml:"clientID"`
	AccessKey              string     `yaml:"accessKey"`
	UploadClientGRPCTarget string     `yaml:"uploadClientGRPCTarget"`
	NumWorkers             int        `yaml:"numWorkers"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := SyncConfig{
		SmbConfig: &SmbConfig{
			User:       "guest",
			Password:   "guest",
			Address:    "10.10.10.204:445",
			ShareName:  "storage",
			DirPath:    "photos",
			Extensions: []string{".jpg", ".jpeg", ".png", ".bmp"},
		},
		SqliteFile:             "files.db",
		ClientID:               "mbp-kkiling-OZON-HXW066MJFG",
		AccessKey:              "1234567",
		UploadClientGRPCTarget: "localhost:8181",
		NumWorkers:             16,
	}

	smb := smbread.NewSmbRead(smbread.Config{
		User:       cfg.SmbConfig.User,
		Password:   cfg.SmbConfig.Password,
		Address:    cfg.SmbConfig.Address,
		ShareName:  cfg.SmbConfig.ShareName,
		DirPath:    cfg.SmbConfig.DirPath,
		Extensions: cfg.SmbConfig.Extensions,
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
		DSN: cfg.SqliteFile,
	})

	if err != nil {
		fmt.Printf("fail sqlitestorage.NewStorage: %v", err)
		return
	}

	// Создание подключения к gRPC серверу
	conn, err := grpc.Dial(cfg.UploadClientGRPCTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		fmt.Printf("failed to connect to gRPC server: %v", err)
		return
	}
	defer conn.Close()

	grpcClient := pbv1.NewSyncPhotosServiceClient(conn)
	//

	uploadClient := syncclient.NewClient(grpcClient, cfg.ClientID, cfg.AccessKey)

	//
	var sync = syncfiles.NewSyncPhotos(smb, storage, uploadClient, syncfiles.Config{
		NumWorkers: cfg.NumWorkers,
	})

	if err := sync.Sync(ctx); err != nil {
		fmt.Printf("fail sync.Sync: %v", err)
		return
	}
}
