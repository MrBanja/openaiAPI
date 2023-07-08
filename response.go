package openaiAPI

import "sync"

type ResponseStream struct {
	dataCh chan string
	errCh  chan error

	mu sync.RWMutex

	isClosedSt bool
}

func newResponseStream() *ResponseStream {
	return &ResponseStream{
		dataCh: make(chan string),
		errCh:  make(chan error),
	}
}

// Data returns a channel that will receive responses from OpenAI.
func (r *ResponseStream) Data() <-chan string {
	return r.dataCh
}

// Error returns a channel that will receive errors from OpenAI or context cancellation.
func (r *ResponseStream) Error() <-chan error {
	return r.errCh
}

func (r *ResponseStream) isClosed() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.isClosedSt
}

func (r *ResponseStream) send(data string) {
	if r.isClosed() {
		return
	}
	r.dataCh <- data
}

func (r *ResponseStream) sendError(err error) {
	if r.isClosed() {
		return
	}
	r.errCh <- err
	r.close()
}

func (r *ResponseStream) close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	close(r.dataCh)
	close(r.errCh)
	r.isClosedSt = true
}
