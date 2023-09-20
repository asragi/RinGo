package test

type TestRandom struct {
	Value float32
}

func (t *TestRandom) Emit() float32 {
	return t.Value
}
