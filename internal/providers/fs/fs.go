package fs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	. "github.com/outscale/frieza/internal/common"
	"github.com/teris-io/cli"
)

const Name = "fs"

const (
	typeFile   = "file"
	typeFolder = "folder"
)

type FileSystem struct {
	Path string
}

func checkConfig(config ProviderConfig) error {
	if len(config["path"]) == 0 {
		return errors.New("path is needed")
	}
	return nil
}

func New(config ProviderConfig, debug bool) (*FileSystem, error) {
	if err := checkConfig(config); err != nil {
		return nil, err
	}
	return &FileSystem{
		Path: config["path"],
	}, nil
}

func Types() []ObjectType {
	object_types := []ObjectType{typeFile, typeFolder}
	return object_types
}

func Cli() (string, cli.Command) {
	return Name, cli.NewCommand(Name, "create new file system profile").
		WithOption(cli.NewOption("path", "folder path"))
}

func (provider *FileSystem) Name() string {
	return Name
}

func (provider *FileSystem) Types() []ObjectType {
	return Types()
}

func (provider *FileSystem) AuthTest(ctx context.Context) error {
	// Will test if we can access the folder
	if _, err := os.ReadDir(provider.Path); err != nil {
		return fmt.Errorf("cannot move to directory: %s", err.Error())
	}
	return nil
}

func (provider *FileSystem) ReadObjects(ctx context.Context, typeName string) ([]Object, error) {
	switch typeName {
	case typeFile:
		return provider.readFiles(ctx)
	case typeFolder:
		return provider.readFolders(ctx)
	}
	return []Object{}, nil
}

func (provider *FileSystem) DeleteObjects(ctx context.Context, typeName string, objects []Object) {
	switch typeName {
	case typeFile:
		provider.deleteFiles(ctx, objects)
	case typeFolder:
		provider.deleteFolders(ctx, objects)
	}
}

func (provider *FileSystem) StringObject(object string, typeName string) string {
	return object
}

func (provider *FileSystem) readFiles(ctx context.Context) ([]Object, error) {
	files := make([]Object, 0)

	if err := os.Chdir(provider.Path); err != nil {
		return nil, fmt.Errorf("chdir: %w", err)
	}

	folderStack := []string{"."}
	for len(folderStack) > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		dirPath := folderStack[len(folderStack)-1]
		folderStack = folderStack[:len(folderStack)-1]
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("read dir: %w", err)
		}
		for _, node := range dir {
			nodePath := path.Join(dirPath, node.Name())
			if node.IsDir() {
				folderStack = append(folderStack, nodePath)
			}
			if node.Type().IsRegular() {
				files = append(files, nodePath)
			}
		}
	}
	return files, nil
}

func (provider *FileSystem) deleteFiles(ctx context.Context, files []Object) {
	for _, relativeFilePath := range files {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled: %v\n", ctx.Err())
			return
		default:
		}

		filePath := path.Join(provider.Path, relativeFilePath)
		log.Printf("Deleting file %s ... ", filePath)
		if err := os.Remove(filePath); err != nil {
			log.Printf("cannot remove file %s\n", err.Error())
			continue
		}
		log.Println("OK")
	}
}

func (provider *FileSystem) readFolders(ctx context.Context) ([]Object, error) {
	folders := make([]Object, 0)

	if err := os.Chdir(provider.Path); err != nil {
		return nil, fmt.Errorf("chdir: %w", err)
	}

	folderStack := []string{"."}
	for len(folderStack) > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		dirPath := folderStack[len(folderStack)-1]
		folderStack = folderStack[:len(folderStack)-1]
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			log.Printf("cannot read directory: %s", err.Error())
			continue
		}
		for _, node := range dir {
			nodePath := path.Join(dirPath, node.Name())
			if node.IsDir() {
				folderStack = append(folderStack, nodePath)
				folders = append(folders, nodePath)
			}
		}
	}
	return folders, nil
}

func (provider *FileSystem) deleteFolders(ctx context.Context, folders []Object) {
	for _, relativeFolderPath := range folders {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled: %v\n", ctx.Err())
			return
		default:
		}

		folderPath := path.Join(provider.Path, relativeFolderPath)
		log.Printf("Deleting folder %s ... ", folderPath)
		if err := os.Remove(folderPath); err != nil {
			log.Printf("cannot remove folder %s\n", err.Error())
			continue
		}
		log.Println("OK")
	}
}
