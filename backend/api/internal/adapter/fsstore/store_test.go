package fsstore

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStore_SaveFile(t *testing.T) {
	ctx := context.Background()
	fileName := fmt.Sprintf("%s.txt", uuid.NewString())
	body := fmt.Sprintf("test text %s", uuid.NewString())
	dir, err := ioutil.TempDir("", uuid.NewString())
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	f := &Store{
		cfg: Config{BaseFilesDir: dir},
	}

	fileKey, err := f.SaveFile(ctx, fileName, []byte(body), "one", "two", "three")
	require.NoError(t, err)
	require.Equal(t, fileKey, filepath.Join("one", "two", "three", fileName))
}

func TestStore_GetFileBody(t *testing.T) {
	ctx := context.Background()
	fileName := fmt.Sprintf("%s.txt", uuid.NewString())
	body := fmt.Sprintf("test text %s", uuid.NewString())
	dir, err := ioutil.TempDir("", uuid.NewString())
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	f := &Store{
		cfg: Config{BaseFilesDir: dir},
	}

	fileKey, err := f.SaveFile(ctx, fileName, []byte(body), "one", "two", "three")
	require.NoError(t, err)

	fileBody, err := f.GetFileBody(ctx, fileKey)
	require.NoError(t, err)
	require.Equal(t, body, string(fileBody))
}

func TestStore_DeleteFile(t *testing.T) {
	ctx := context.Background()
	fileName := fmt.Sprintf("%s.txt", uuid.NewString())
	body := fmt.Sprintf("test text %s", uuid.NewString())
	dir, err := ioutil.TempDir("", uuid.NewString())
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	f := &Store{
		cfg: Config{BaseFilesDir: dir},
	}

	fileKey, err := f.SaveFile(ctx, fileName, []byte(body), "one", "two", "three")
	require.NoError(t, err)

	err = f.DeleteFile(ctx, fileKey)
	require.NoError(t, err)

	_, err = f.GetFileBody(ctx, fileKey)
	require.Error(t, err)
}
