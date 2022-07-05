package gcputil

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
)

type Project struct {
	id   string
	name string
}

func NewProject(projectId string) (*Project, error) {
	project := Project{
		id: projectId,
	}

	ctx := context.Background()
	crm, err := cloudresourcemanager.NewService(ctx)

	if err != nil {
		return nil, errors.Wrapf(err, "could not create cloudresourcemanager service")
	}

	projectService := cloudresourcemanager.NewProjectsService(crm)
	gcpProject, err := projectService.Get(projectId).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "could not get project name from project id %s", projectId)
	}

	project.name = gcpProject.Name

	return &project, nil
}

func (a *Project) ID() string {
	return a.id
}

func (a *Project) Name() string {
	return a.name
}
