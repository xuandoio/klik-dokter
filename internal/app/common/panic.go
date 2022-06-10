package common

import "sync"

type PanicGroup struct {
	wg     sync.WaitGroup
	errOne sync.Once
	err    error
}

// Wait /**
func (g *PanicGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

// Go /**
func (g *PanicGroup) Go(f func()) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					g.errOne.Do(func() {
						g.err = err
					})
				} else {
					panic(r)
				}
			}
		}()

		f()
	}()
}
