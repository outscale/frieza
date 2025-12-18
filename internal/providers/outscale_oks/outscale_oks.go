package outscale_oks

import (
	"context"
	"log"

	. "github.com/outscale/frieza/internal/common"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	oscutils "github.com/outscale/osc-sdk-go/v3/pkg/utils"
	"github.com/teris-io/cli"
)

const (
	Name = "outscale_oks"

	typeProject = "project"
	typeCluster = "cluster"
)

type OutscaleOKS struct {
	client *oks.Client
}

func (provider *OutscaleOKS) StringObject(object string, typeName string) string {
	return object
}

func New(config ProviderConfig, debug bool) (*OutscaleOKS, error) {
	profileName := config["profile"]
	profilePath := config["path"]
	profile, err := profile.NewProfileFromStandardConfiguration(profileName, profilePath)
	if err != nil {
		return nil, err
	}

	if ak, ok := config["ak"]; ok {
		profile.AccessKey = ak
	}

	if sk, ok := config["sk"]; ok {
		profile.SecretKey = sk
	}

	if region, ok := config["region"]; ok {
		profile.Region = region
	}

	client, err := oks.NewClient(profile, oscutils.WithUseragent("frieza/"+FullVersion()))
	if err != nil {
		return nil, err
	}

	return &OutscaleOKS{
		client: client,
	}, nil
}

func Types() []ObjectType {
	return []ObjectType{
		typeProject,
		typeCluster,
	}
}

func Cli() (string, cli.Command) {
	return Name, cli.NewCommand(Name, "create new Outscale OKS profile").
		WithOption(cli.NewOption("region", "Outscale region (e.g. eu-west-2)")).
		WithOption(cli.NewOption("ak", "access key")).
		WithOption(cli.NewOption("sk", "secret key"))
}

func (provider *OutscaleOKS) Name() string {
	return Name
}

func (provider *OutscaleOKS) Types() []ObjectType {
	return Types()
}

func (provider *OutscaleOKS) AuthTest() error {
	// TODO
	return nil
}

func (provider *OutscaleOKS) ReadObjects(typeName string) ([]Object, error) {
	switch typeName {
	case typeProject:
		return provider.readProject()
	case typeCluster:
		return provider.readCluster()
	}

	return []Object{}, nil
}

func (provider *OutscaleOKS) readCluster() ([]Object, error) {
	clusters := make([]Object, 0)
	read, err := provider.client.ListAllClusters(context.Background(), &oks.ListAllClustersParams{})
	if err != nil {
		return clusters, err
	}

	for _, cluster := range read.Clusters {
		clusters = append(clusters, cluster.Id)
	}

	return clusters, nil
}

func (provider *OutscaleOKS) readProject() ([]Object, error) {
	projectcs := make([]Object, 0)
	read, err := provider.client.ListProjects(context.Background(), &oks.ListProjectsParams{})
	if err != nil {
		return projectcs, err
	}

	for _, project := range read.Projects {
		projectcs = append(projectcs, project.Id)
	}

	return projectcs, nil
}

func (provider *OutscaleOKS) DeleteObjects(typeName string, objects []Object) {
	switch typeName {
	case typeProject:
		provider.deleteProject(objects)
	case typeCluster:
		provider.deleteCluster(objects)
	}
}

func (provider *OutscaleOKS) deleteCluster(objects []Object) {
	if len(objects) == 0 {
		return
	}

	ctx := context.Background()
	for _, clusterID := range objects {
		log.Printf("Deleting cluster %s... ", clusterID)

		_, err := provider.client.DeleteCluster(ctx, clusterID)
		if err != nil {
			log.Printf("Error while deleting cluster: %v\n", err)
		} else {
			log.Println("OK")
		}
	}
}

func (provider *OutscaleOKS) deleteProject(objects []Object) {
	if len(objects) == 0 {
		return
	}

	ctx := context.Background()
	for _, projectID := range objects {
		log.Printf("Deleting project %s... ", projectID)

		_, err := provider.client.DeleteProject(ctx, projectID)
		if err != nil {
			log.Printf("Error while deleting project: %v\n", err)
		} else {
			log.Println("OK")
		}
	}
}
