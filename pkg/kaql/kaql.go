package kaql

type KaQL struct {
	Query Node
}

func New(src string) (*KaQL, error) {
	node, err := NewParser().Parse(src)
	if err != nil {
		return nil, err
	}

	return &KaQL{Query: node}, nil
}

func (q *KaQL) IsEmpty() bool {
	return q.Query == nil
}

func (q *KaQL) Match(obj interface{}) bool {
	if q.IsEmpty() {
		return true
	}

	return true
}
