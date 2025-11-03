package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	. "github.com/outscale/frieza/internal/common"
)

type Destroyer struct {
	Targets []DestroyerTarget `json:"targets"`
}

type DestroyerTarget struct {
	profile     *Profile          `json:"-"`
	JsonProfile *DestroyerProfile `json:"profile"`
	provider    *Provider         `json:"-"`
	Objects     *Objects          `json:"objects"`
}

type DestroyerProfile struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

func NewDestroyer() *Destroyer {
	var destroyer Destroyer
	return &destroyer
}

func (destroyer *Destroyer) add(profile *Profile, provider *Provider, objectsToDelete *Objects) {
	target := DestroyerTarget{
		profile: profile,
		JsonProfile: &DestroyerProfile{
			Name:     profile.Name,
			Provider: profile.Provider,
		},
		provider: provider,
		Objects:  objectsToDelete,
	}
	destroyer.Targets = append(destroyer.Targets, target)
}

func (destroyer *Destroyer) print(json bool) {
	if json {
		destroyer.print_json()
	} else {
		destroyer.print_human()
	}
}

func (destroyer *Destroyer) print_human() {
	count := len(destroyer.Targets)
	totalObjectCount := 0
	for i := 0; i < count; i++ {
		target := destroyer.Targets[i]
		log.Printf(
			"Objects to delete in profile %s (%s):\n",
			target.profile.Name,
			(*target.provider).Name(),
		)
		objectsCount := ObjectsCount(target.Objects)
		if objectsCount == 0 {
			log.Println("* no object *")
		}
		totalObjectCount += objectsCount
		log.Print(ObjectsPrint(target.provider, target.Objects))
	}
	if totalObjectCount == 0 {
		log.Println("\nNothing to delete")
	}
}

func (destroyer *Destroyer) print_json() {
	json_bytes, err := json.MarshalIndent(destroyer, "", "  ")
	if err != nil {
		cliFatalf(true, "Cannot serialize to json: %s", err.Error())
	}
	log.Print(string(json_bytes))
}

func (destroyer *Destroyer) run() {
	var objects []*Objects
	for i := range destroyer.Targets {
		objects = append(objects, destroyer.Targets[i].Objects)
	}
	var hasObjectsLeft []bool
	for range objects {
		hasObjectsLeft = append(hasObjectsLeft, true)
	}
	for {
		var objectsCount []int
		var totalObjectCount int
		for i := range objects {
			objectsCount = append(objectsCount, 0)
			if !hasObjectsLeft[i] {
				continue
			}
			count := ObjectsCount(objects[i])
			objectsCount[i] = count
			if count == 0 {
				hasObjectsLeft[i] = false
				continue
			}
			totalObjectCount += count
		}
		if totalObjectCount == 0 {
			return
		}
		for i, target := range destroyer.Targets {
			if !hasObjectsLeft[i] {
				continue
			}
			DeleteObjects(target.provider, *objects[i])
			time.Sleep(100 * time.Millisecond)
		}
		for i, target := range destroyer.Targets {
			diff := NewDiff()
			remaining, err := ReadNonEmptyObjects(target.provider, *objects[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
			}
			diff.Build(&remaining, objects[i])
			objects[i] = &diff.Retained
		}
		time.Sleep(time.Second)
	}
}

func confirmAction(message *string, autoApprove bool) bool {
	if autoApprove {
		return true
	}
	log.Printf("\n%s\n", *message)
	log.Printf("  There is no undo. Only 'yes' will be accepted to confirm.\n\n")
	log.Printf("  Enter a value: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ReplaceAll(response, "\n", "")
	response = strings.ReplaceAll(response, "\r", "")
	return response == "yes"
}
