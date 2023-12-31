package stream

/*
Implements:
  - Multiplexer
*/
type multiplexer[T any] struct {
	DefaultConsumer[T]
	outputs []ChanneledInput[T]
}

/*
NewMultiplexer is a constructor of the multiplexer.

Type parameters:
  - T - type of the consumed and produced values.

Parameters:
  - capacity - size of the channel buffer,
  - branches - number of the output streams.

Returns:
  - pointer to the new multiplexer.
*/
func NewMultiplexer[T any](capacity int, branches int) Multiplexer[T] {
	ego := &multiplexer[T]{}
	ego.outputs = make([]ChanneledInput[T], branches)
	for i := 0; i < branches; i++ {
		ego.outputs[i] = NewChanneledInput[T](capacity)
	}
	return ego
}

/*
Consumes the data from the source Producer and pushes them to the result streams.
Every value is pushed to all streams.
It runs asynchronously.
*/
func (ego *multiplexer[T]) pipeData() {

	for _, output := range ego.outputs {
		defer output.Close()
	}

	for {
		value, valid, err := ego.Consume()
		if err != nil || !valid {
			return
		}
		for _, output := range ego.outputs {
			output.Write(value)
		}
	}

}

func (ego *multiplexer[T]) Out(index int) Producer[T] {
	return ego.outputs[index]
}

func (ego *multiplexer[T]) SetSource(s Producer[T]) error {
	if err := ego.DefaultConsumer.SetSource(s); err != nil {
		return err
	}
	go ego.pipeData()
	return nil
}
