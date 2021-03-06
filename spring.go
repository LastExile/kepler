package kepler

import (
	"context"
)

// Spring of outgong messages, act as producer
type Spring interface {
	Out(ctx context.Context, out chan Message) <-chan Message
	LinkTo(name string, target Sink, cond RouteCondition) (closer func())
}

func (s *springImpl) Out(ctx context.Context, o chan Message) <-chan Message {

	go func() {
		s.action(ctx, o)
	}()

	return o
}

func (s *springImpl) LinkTo(name string, sink Sink, cond RouteCondition) (closer func()) {
	route := s.router.AddRoute(name, cond)

	//pass linked conext to Sink
	inCtx, inClose := context.WithCancel(route.Ctx())
	sink.In(inCtx, s.Out(route.Ctx(), route.Buff()))

	//send signal to Sink and close original route
	return func() { inClose(); route.Close() }
}

// SpringFunction out generator function
type SpringFunction func(ctx context.Context, out chan<- Message)

// UnmarshalFunction used to unmarshal spring input to
type UnmarshalFunction func(in []byte) (Message, error)

// NewSpring creates new Spring
func NewSpring(action SpringFunction) Spring {
	return &springImpl{action: action, router: NewRouter(false)}
}

type springImpl struct {
	action SpringFunction
	router Router
}
