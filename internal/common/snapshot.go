package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Snapshot struct {
	Version int            `json:"version"`
	Name    string         `json:"name"`
	Date    string         `json:"date"`
	Data    []SnapshotData `json:"data"`
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

func (diff *Diff) String() string {
	out := fmt.Sprintf("  - retained:\n")
	out += ObjectsPrint(&diff.Retained)
	out += fmt.Sprintf("  - created:\n")
	out += ObjectsPrint(&diff.Created)
	out += fmt.Sprintf("  - deleted:\n")
	out += ObjectsPrint(&diff.Deleted)
	return out
}

func ObjectsCount(objects *Objects) int {
	count := 0
	for _, objectIds := range *objects {
		count += len(objectIds)
	}
	return count
}

func ObjectsPrint(objects *Objects) string {
	out := ""
	for objectType, objectIds := range *objects {
		if len(objectIds) == 0 {
			continue
		}
		out += fmt.Sprintf("    - %s:\n", objectType)
		for _, objectId := range objectIds {
			out += fmt.Sprintf("      - %s\n", objectId)
		}
	}
	return out
}

func DefaultSnapshotFolderPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".frieza", "snapshots"), nil
}

func (snapshot *Snapshot) Write() error {
	snapshotFolderPath, err := DefaultSnapshotFolderPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(snapshotFolderPath, os.ModePerm); err != nil {
		return err
	}
	snapshotPath := path.Join(snapshotFolderPath, snapshot.Name+".json")
	json_bytes, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(snapshotPath, json_bytes, 0700); err != nil {
		return err
	}
	return nil
}

func SnapshotLoad(name string) (*Snapshot, error) {
	snapshotFolderPath, err := DefaultSnapshotFolderPath()
	if err != nil {
		return nil, err
	}
	snapshotPath := path.Join(snapshotFolderPath, name+".json")
	snapshot_json, err := ioutil.ReadFile(snapshotPath)
	if err != nil {
		return nil, err
	}
	var snapshot Snapshot
	if err := json.Unmarshal(snapshot_json, &snapshot); err != nil {
		return nil, err
	}
	if snapshot.Version > SnapshotVersion() {
		return nil, errors.New("snapshot version not supported, please upgrade frieza")
	}
	return &snapshot, nil
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
	snapshotFolderPath, err := DefaultSnapshotFolderPath()
	if err != nil {
		return err
	}
	snapshotPath := path.Join(snapshotFolderPath, snapshot.Name+".json")
	if err = os.Remove(snapshotPath); err != nil {
		return err
	}
	return nil
}
