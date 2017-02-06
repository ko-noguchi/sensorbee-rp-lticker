package plugin

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"github.com/ktanaka89/sensorbee-plugin-psutil"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("psutil", bql.SourceCreatorFunc(sensorbee_plugin_psutil.NewIntervalSource))
}
