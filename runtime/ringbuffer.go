package runtime

type RingBuffer struct {
	buf         []byte
	start, size int
}

func New(size int) *RingBuffer {
	return &RingBuffer{make([]byte, size), 0, 0}
}

func (r *RingBuffer) Write(b []byte) {
	for len(b) > 0 {
		start := (r.start + r.size) % len(r.buf)
		n := copy(r.buf[start:], b)
		b = b[n:] //golang就是要好好运用切片

		if r.size >= len(r.buf) {
			if n <= len(r.buf) {
				r.start += n
				if r.start >= len(r.buf) {
					r.start = 0
				}
			} else {
				r.start = 0
			}
		}
		r.size += n
		// Size can't exceed the capacity
		if r.size > cap(r.buf) {
			r.size = cap(r.buf)
		}
	}
}

func (r *RingBuffer) Read(b []byte) int {
	read := 0
	size := r.size
	start := r.start
	for len(b) > 0 && size > 0 {
		end := start + size
		if end > len(r.buf) {
			end = len(r.buf)
		}
		n := copy(b, r.buf[start:end])
		size -= n
		read += n
		b = b[n:]

		start = (start + n) % len(r.buf)
	}
	return read
}

func (r *RingBuffer) Size() int {
	return r.size
}
