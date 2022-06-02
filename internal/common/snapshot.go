package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Objects = map[ObjectType][]Object

type Snapshot struct {
	Version int            `json:"version"`
	Name    string         `json:"name"`
	Date    string         `json:"date"`
	Data    []SnapshotData `json:"data"`
	Config  *Config        `json:"-"`
}

type SnapshotData struct {
	Profile string  `json:"profile"`
	Objects Objects `json:"objects"`
}

type Diff struct {
	Retained Objects
	Created  Objects
	Deleted  Objects
}

func SnapshotVersion() int {
	return 0
}

func ReadObjects(provider *Provider) Objects {
	objects := make(Objects)
	for _, typeName := range (*provider).Types() {
		objects[typeName] = (*provider).ReadObjects(typeName)
	}
	return objects
}

func FiltersObjects(objects *Objects, onlyResources []ObjectType) Objects {
	filteredObjects := make(Objects)
	for _, typeName := range onlyResources {
		filteredObjects[typeName] = (*objects)[typeName]
	}
	return filteredObjects
}

func ReadNonEmptyObjects(provider *Provider, nonEmpy Objects) Objects {
	objects := make(Objects)
	for _, typeName := range (*provider).Types() {
		if len(nonEmpy[typeName]) == 0 {
			continue
		}
		objects[typeName] = (*provider).ReadObjects(typeName)
	}
	return objects
}

func DeleteObjects(provider *Provider, objects Objects) {
	for _, typeName := range (*provider).Types() {
		objectList := objects[typeName]
		if len(objectList) == 0 {
			continue
		}
		(*provider).DeleteObjects(typeName, objectList)
	}
}

func NewDiff() *Diff {
	return &Diff{
		Retained: make(Objects),
		Created:  make(Objects),
		Deleted:  make(Objects),
	}
}

func objects2Map(objectIds []Object) map[string]bool {
	out := make(map[string]bool)
	for _, objectId := range objectIds {
		out[objectId] = true
	}
	return out
}

func (diff *Diff) Build(a *Objects, b *Objects) {
	allTypes := make(map[string]bool)
	for objectTypeA := range *a {
		allTypes[objectTypeA] = true
	}
	for objectTypeB := range *b {
		allTypes[objectTypeB] = true
	}

	for objectType := range allTypes {
		aFlat := objects2Map((*a)[objectType])
		bFlat := objects2Map((*b)[objectType])
		for idA := range aFlat {
			if bFlat[idA] {
				diff.Retained[objectType] = append(diff.Retained[objectType], idA)
			} else {
				diff.Deleted[objectType] = append(diff.Deleted[objectType], idA)
			}
		}
		for idB := range bFlat {
			if !aFlat[idB] {
				diff.Created[objectType] = append(diff.Created[objectType], idB)
			}
		}
	}
}

func ObjectsCount(objects *Objects) int {
	count := 0
	for _, objectIds := range *objects {
		count += len(objectIds)
	}
	return count
}

func ObjectsPrint(provider *Provider, objects *Objects) string {
	out := ""
	for objectType, objectIds := range *objects {

		if len(objectIds) == 0 {
			continue
		}
		out += fmt.Sprintf("%s:\n", objectType)
		for _, objectId := range objectIds {
			out += fmt.Sprintf("  - %s\n", (*provider).StringObject(objectId, objectType))
		}
	}
	return out
}

func (snapshot *Snapshot) Write() error {
	if err := os.MkdirAll(snapshot.Config.SnapshotFolderPath, os.ModePerm); err != nil {
		return err
	}
	json_bytes, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(snapshot.Path(), json_bytes, 0700); err != nil {
		return err
	}
	return nil
}

func (snapshot *Snapshot) Path() string {
	return path.Join(snapshot.Config.SnapshotFolderPath, snapshot.Name+".json")
}

func SnapshotLoad(name string, config *Config) (*Snapshot, error) {
	snapshot := &Snapshot{
		Name:   name,
		Config: config,
	}
	snapshot_json, err := ioutil.ReadFile(snapshot.Path())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(snapshot_json, &snapshot); err != nil {
		return nil, err
	}
	if snapshot.Version > SnapshotVersion() {
		return nil, errors.New("snapshot version not supported, please upgrade frieza")
	}
	return snapshot, nil
}

func (snapshot Snapshot) String() string {
	out := fmt.Sprintf("name: %v\n", snapshot.Name)
	out += fmt.Sprintf("date: %v\n", snapshot.Date)
	out += fmt.Sprintf("profiles:\n")
	for _, data := range snapshot.Data {
		out += fmt.Sprintf("  - %v:\n", data.Profile)
		for objectType, objects := range data.Objects {
			out += fmt.Sprintf("    - %s: %d\n", objectType, len(objects))
		}
	}
	return out
}

func (snapshot Snapshot) Delete() error {
	return os.Remove(snapshot.Path())
}
