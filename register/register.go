package register

import (
	"context"
	"fmt"
	"github/inoth/ino-gateway/components"
	"github/inoth/ino-gateway/service"
	"os"
	"sync"
)

var initOnce sync.Once

type Register struct {
	cmps []components.Components
	svcs []service.Service
}

func NewRegister(cmps ...components.Components) *Register {
	return &Register{
		cmps: cmps,
	}
}

func (reg *Register) Init() *Register {
	initOnce.Do(func() {
		for _, cmp := range reg.cmps {
			must(cmp.Init())
		}
	})
	return reg
}

func (reg *Register) SubStart(svcs ...service.Service) *Register {
	reg.svcs = svcs
	for _, svc := range reg.svcs {
		go func(ctx context.Context, service service.Service) {
			defer func() {
				if exception := recover(); exception != nil {
					if err, ok := exception.(error); ok {
						fmt.Printf("%v\n", err)
					} else {
						panic(exception)
					}
					os.Exit(1)
				}
			}()
			// run sub service
			must(service.Start())
		}(context.Background(), svc)
	}
	return reg
}

func (reg *Register) Start(svcs service.Service) {
	must(svcs.Start())
}

func (reg *Register) Stop() {
	for _, svcs := range reg.svcs {
		svcs.Stop()
	}
}

func must(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
