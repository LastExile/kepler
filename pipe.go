package kepler

import (
	"context"
	"log"
)

type Pipe interface {
	Spring
	Sink
}

type pipeImpl struct {
	name      string
	action    PipeFunction
	router    Router
	broadcast bool
}

// PipeFunction defines transform func that will be performed within block
type PipeFunction func(in Message) Message

// Out outgoing channel
func (p *pipeImpl) Out(ctx context.Context, buf chan Message) <-chan Message {
	return buf
}

// In incomming channel
func (p *pipeImpl) In(ctx context.Context, input <-chan Message) {
	go func() {
		for {
			select {
			case msg := <-input:
				if msg != nil {
					m := p.action(msg)
					if m != nil {
						p.router.Send(m)
					}
				}
			case <-ctx.Done():
				log.Println("In Done")
				//TODO: propogate to attached
				p.router.Close()
				return
			}
		}
	}()
}

// Name of this pipe
func (p *pipeImpl) Name() string {
	return p.name
}

// LinkTo add new conditional link
func (p *pipeImpl) LinkTo(sink Sink, cond RouteCondition) (closer func()) {
	route := p.router.AddRoute(sink.Name(), cond)

	//inCtx, inClose := context.WithCancel(route.Ctx())

	//pass to target sink route context
	sink.In(route.Ctx(), p.Out(nil, route.Buff()))

	return func() { route.Close() }
}

// NewPipe creates new instance of pipe with defined transform action
func NewPipe(name string, action PipeFunction) Pipe {
	return &pipeImpl{name: name, action: action, router: NewRouter(false)}
}

func NewBroadcastPipe(name string, action PipeFunction) Pipe {
	return &pipeImpl{name: name, action: action, router: NewRouter(true)}
}
