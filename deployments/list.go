package deployments

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	deployments, err := c.DeploymentList(app)
	if err != nil {
		return errgo.Mask(err)
	}

	if len(deployments) == 0 {

	} else {
		t := tablewriter.NewWriter(os.Stdout)
		t.SetHeader([]string{"ID", "Date", "User", "Git Ref", "Status"})

		for _, deployment := range deployments {
			t.Append([]string{deployment.ID,
				deployment.CreatedAt.Format("2006/01/02 15:04:05"),
				deployment.User.Username,
				deployment.GitRef,
				string(deployment.Status),
			})
		}
		t.Render()

	}

	return nil
}
