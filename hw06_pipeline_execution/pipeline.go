package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func doneStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case prev, ok := <-in:
				if !ok {
					return
				}
				out <- prev
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipe := in
	for _, stage := range stages {
		if stage != nil {
			pipe = stage(doneStage(pipe, done))
		}
	}
	return pipe
}
