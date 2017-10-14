package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
)

type RunConf struct {
	Port     string
	Sess     *session.Session
	KmsKeyID string
}

func New() *RunConf {
	ret := new(RunConf)

	port := getenv.String("PNZR_PORT", "8080")
	ret.Port = fmt.Sprintf(":%s", port)

	ret.Sess = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 getenv.String("AWS_PROFILE", "default"),
	}))

	ret.KmsKeyID = getenv.String("KMS_KEY_ID")

	return ret
}
