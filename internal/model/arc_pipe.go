package model

// arcPipe is closed before the caller receives the arc length; a follow-up
// tag write after Close panics.
type arcPipe struct {
	closed bool
	tags   map[string]float64
}

func (p *arcPipe) Close() {
	p.closed = true
	p.tags = nil
}

func (p *arcPipe) tagLength(name string, length float64) {
	p.tags[name] = length
}

func sealArcPipe(length float64) {
	p := &arcPipe{tags: map[string]float64{}}
	defer p.Close()
	p.Close()
	p.tagLength("arc", length)
}
