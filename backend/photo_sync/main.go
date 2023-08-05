package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/hirochachacha/go-smb2"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"strings"
)

var imageExtensions = []string{".jpg", ".png", ".gif", ".bmp", ".jpeg"}

func readDir(fs *smb2.Share, dirPath string, files chan<- string) error {
	dir, err := fs.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		newPath := dirPath + "/" + name
		if entry.IsDir() {
			err = readDir(fs, newPath, files)
			if err != nil {
				return err
			}
		} else {
			for _, ext := range imageExtensions {
				if strings.HasSuffix(strings.ToLower(name), ext) {
					files <- newPath
				}
			}
		}
	}

	return nil
}

func main() {
	ctx := context.Background()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "guest",
			Password: "guest",
		},
	}

	tconn, err := net.Dial("tcp", "nas.lan"+":445")
	if err != nil {
		panic(err)
	}

	conn, err := d.DialContext(ctx, tconn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Logoff()

	s := conn.WithContext(ctx)
	fs, err := s.Mount(`storage`)
	if err != nil {
		log.Fatal(err)
	}

	defer fs.Umount()

	filesChan := make(chan string)

	go func() {
		defer close(filesChan)
		err = readDir(fs, "photos\\Фото\\13.01.18", filesChan)
		if err != nil {
			log.Fatal(err)
		}
	}()

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("Find files"),
	)

	filesAll := make([]string, 0)
	photos := make(map[string][]string, 0)
	for filePath := range filesChan {
		filesAll = append(filesAll, filePath)
		bar.Add(1)
	}
	bar.Finish()

	// Создание подключения к gRPC серверу
	grpcConn, err := grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer grpcConn.Close()

	// Создание клиента
	client := pbv1.NewPhotosServiceClient(grpcConn)

	bar = progressbar.NewOptions(len(filesAll),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("Upload files"),
	)
	for _, filePath := range filesAll {
		file, err := fs.Open(filePath)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}

		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			log.Fatalf("Failed to copy content to hasher: %v", err)
		}

		hashString := fmt.Sprintf("%x", hash.Sum(nil))
		bar.Add(1) // increment the progressbar after processing each element
		file.Close()

		photos[hashString] = append(photos[hashString], filePath)

		response, err := client.CheckHashPhoto(ctx, &pbv1.CheckHashPhotoRequest{Hash: "123"})
		if err != nil {
			log.Fatalf("Failed CheckHashPhotoRequest: %v", err)
		}

		if response.AlreadyUploaded {

		} else {

		}
	}

	fmt.Printf("Total: %d", len(filesAll))
	fmt.Printf("Total unique: %d", len(photos))
}
