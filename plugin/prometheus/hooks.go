package prometheus

import (
	"github.com/xfp-881643/gmqtt/server"
)

func (p *Prometheus) HookWrapper() server.HookWrapper {
	return server.HookWrapper{}
}
