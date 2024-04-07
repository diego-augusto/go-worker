package worker

type optFunc func(*pool)

func WithWorkers(workers ...Worker) optFunc {
	return func(w *pool) {
		w.workers = workers
	}
}

func WithExecuters(executers ...Executer) optFunc {
	return func(w *pool) {
		w.executers = executers
	}
}

func WithDefaultExecuter() optFunc {
	return func(w *pool) {
		w.executers = append(w.executers, NewDefaultExecuter())
	}
}
