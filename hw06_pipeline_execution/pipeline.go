package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type (
	Stage          func(in In) (out Out)
	StageAbortable func(in, done In) Out
)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = makeAbortableStage(stage)(in, done)
	}
	return in
}

func makeAbortableStage(original Stage) StageAbortable {
	return func(in, done In) Out {
		in = original(in)
		out := make(Bi)

		go func() {
			defer close(out)

			for {
				select {
				case x, ok := <-in:
					if !ok {
						return
					}
					out <- x
				case <-done:
					return
				}
			}
		}()

		return out
	}
}
