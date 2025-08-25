package fs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	. "github.com/outscale/frieza/internal/common"
	"github.com/teris-io/cli"
)

const Name = "fs"

const typeFile = "file"
const typeFolder = "folder"

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

func (provider *FileSystem) AuthTest() error {
	// Will test if we can access the folder
	if _, err := os.ReadDir(provider.Path); err != nil {
		return fmt.Errorf("cannot move to directory: %s", err.Error())
	}
	return nil
}

func (provider *FileSystem) ReadObjects(typeName string) []Object {
	switch typeName {
	case typeFile:
		return provider.readFiles()
	case typeFolder:
		return provider.readFolders()
	}
	return []Object{}
}

func (provider *FileSystem) DeleteObjects(typeName string, objects []Object) {
	switch typeName {
	case typeFile:
		provider.deleteFiles(objects)
	case typeFolder:
		provider.deleteFolders(objects)
	}
}

func (provider *FileSystem) StringObject(object string, typeName string) string {
	return object
}

func (provider *FileSystem) readFiles() []Object {
	files := make([]Object, 0)

	if err := os.Chdir(provider.Path); err != nil {
		log.Printf("cannot move to directory: %s", err.Error())
		return files
	}

	folderStack := []string{"."}
	for len(folderStack) > 0 {
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
			}
			if node.Type().IsRegular() {
				files = append(files, nodePath)
			}
		}
	}
	return files
}

func (provider *FileSystem) deleteFiles(files []Object) {
	for _, relativeFilePath := range files {
		filePath := path.Join(provider.Path, relativeFilePath)
		log.Printf("Deleting file %s ... ", filePath)
		if err := os.Remove(filePath); err != nil {
			log.Printf("cannot remove file %s\n", err.Error())
			continue
		}
		log.Println("OK")
	}
}

func (provider *FileSystem) readFolders() []Object {
	folders := make([]Object, 0)

	if err := os.Chdir(provider.Path); err != nil {
		log.Printf("cannot move to directory: %s", err.Error())
		return folders
	}

	folderStack := []string{"."}
	for len(folderStack) > 0 {
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
	return folders
}

func (provider *FileSystem) deleteFolders(folders []Object) {
	for _, relativeFolderPath := range folders {
		folderPath := path.Join(provider.Path, relativeFolderPath)
		log.Printf("Deleting folder %s ... ", folderPath)
		if err := os.Remove(folderPath); err != nil {
			log.Printf("cannot remove folder %s\n", err.Error())
			continue
		}
		log.Println("OK")
	}
}
