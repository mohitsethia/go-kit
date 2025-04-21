package pkg

type Options struct {
	withUnInterruptedContext bool
}

type Option func(*Options)

func WithUnInterruptedContext() Option {
	return func(o *Options) {
		o.withUnInterruptedContext = true
	}
}
