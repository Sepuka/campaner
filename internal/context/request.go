package context

type Request struct {
	kind string
}

func (o *Request) GetKind() string {
	return o.kind
}
