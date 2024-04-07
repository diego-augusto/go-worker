package worker

type optFunc func(*worker)

func WithDoers(doers ...Doer) optFunc {
	return func(w *worker) {
		w.doers = doers
	}
}

func WithExecuters(executers ...Executer) optFunc {
	return func(w *worker) {
		w.executers = executers
	}
}

func WithDefaultExecuter() optFunc {
	return func(w *worker) {
		w.executers = append(w.executers, NewDefaultExecuter())
	}
}
