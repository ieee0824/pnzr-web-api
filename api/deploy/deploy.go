package deploy

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ieee0824/errors"
	"github.com/ieee0824/pnzr-web-api/lib/config"
	"github.com/jobtalk/pnzr/api"
	"github.com/jobtalk/pnzr/lib/setting/v1"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)



func parseDockerImage(image string) (url, tag string, err error) {
	r := strings.Split(image, ":")
	if 3 <= len(r) {
		return "", "", errors.New("docker image parse error")
	}
	if len(r) == 2 {
		return r[0], r[1], nil
	}
	return r[0], "", nil
}

type DeployRequest struct {
	Base []byte   `json:"base"`
	Vars [][]byte `json:"vars"`
	Tag  string   `json:"tag"`
}

func Deploy(ctx *gin.Context) {
	cfg := config.New()
	loader := v1.NewLoader(cfg.Sess, &cfg.KmsKeyID)
	now := time.Now().UnixNano()
	tempDir := fmt.Sprintf("%s/%d", os.TempDir(), now)
	varsDir := fmt.Sprintf("%s/vars", tempDir)
	req := DeployRequest{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		panic(err)
	}

	basePath := fmt.Sprintf("%s/base.json", tempDir)

	if err := ioutil.WriteFile(basePath, req.Base, 0644); err != nil {
		panic(err)
	}

	for i, data := range req.Vars {
		path := fmt.Sprintf("%s/%04d.json", varsDir, i)
		if err := ioutil.WriteFile(path, data, 0644); err != nil {
			panic(err)
		}
	}

	s, err := loader.Load(basePath, varsDir, "")
	if err != nil {
		panic(err)
	}

	for i, containerDefinition := range s.TaskDefinition.ContainerDefinitions {
		imageName, tag, err := parseDockerImage(*containerDefinition.Image)
		if err != nil {
			panic(err)
		}

		if tag == "$tag" {
			image := imageName + ":" + req.Tag
			s.TaskDefinition.ContainerDefinitions[i].Image = &image
		} else if tag == "" {
			image := imageName + ":" + "latest"
			s.TaskDefinition.ContainerDefinitions[i].Image = &image
		}
	}

	result, err := api.Deploy(cfg.Sess, s)
	if err != nil {
		panic(err)
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(
		http.StatusOK,
		result,
	)
}
