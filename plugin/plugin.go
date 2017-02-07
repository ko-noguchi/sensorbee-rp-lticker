package plugin

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"github.com/ko-noguchi/sensorbee-rp-lticker"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("rp-lticker", bql.SourceCreatorFunc(sensorbee_rp_lticker.NewIntervalSource))
}
