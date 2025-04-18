package pkg

type Options struct {
	withUnInterruptedContext bool
}

func WithUnInterruptedContext() Options {
	return Options{
		withUnInterruptedContext: true,
	}
}
