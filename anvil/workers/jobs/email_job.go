package jobs

import (
	"github.com/sampiiiii-dev/anvil_server/anvil/common"
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
	"github.com/sampiiiii-dev/anvil_server/anvil/logs"
)

type EmailJob struct {
	Email string
}

func (e *EmailJob) Execute() error {
	s := logs.HireScribe()
	c := config.GetConfigInstance(s)
	smtpClient := common.GetSMTPClient(c)
	return smtpClient.SendEmail(e.Email, "Subject here", "<h1>Hello</h1>", true)
}
