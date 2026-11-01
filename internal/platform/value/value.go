// Package value provides copies of optional scalar values for in-memory repositories.
package value

func Copy[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}
