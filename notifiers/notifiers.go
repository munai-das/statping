package notifiers

import (
	"bytes"
	"github.com/statping/statping/types/core"
	"github.com/statping/statping/types/failures"
	"github.com/statping/statping/types/services"
	"github.com/statping/statping/utils"
	"html/template"
	"time"
)

var log = utils.Log.WithField("type", "notifier")

type replacer struct {
	Core    core.Core
	Service services.Service
	Failure failures.Failure
	Custom  map[string]string
}

func InitNotifiers() {
	Add(
		slacker,
		Command,
		Discorder,
		email,
		LineNotify,
		Telegram,
		Twilio,
		Webhook,
		Mobile,
		Pushover,
		statpingMailer,
		Gotify,
	)
}

func ReplaceTemplate(tmpl string, data replacer) string {
	buf := new(bytes.Buffer)
	tmp, err := template.New("replacement").Parse(tmpl)
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	err = tmp.Execute(buf, data)
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	return buf.String()
}

func Add(notifs ...services.ServiceNotifier) {
	for _, n := range notifs {
		services.AddNotifier(n)
		if err := n.Select().Create(); err != nil {
			log.Error(err)
		}
	}
}

func ReplaceVars(input string, s services.Service, f failures.Failure) string {
	return ReplaceTemplate(input, replacer{Service: s, Failure: f, Core: *core.App})
}

var exampleFailure = &failures.Failure{
	Id:        1,
	Issue:     "HTTP returned a 500 status code",
	ErrorCode: 500,
	Service:   1,
	PingTime:  43203,
	CreatedAt: utils.Now().Add(-10 * time.Minute),
}
