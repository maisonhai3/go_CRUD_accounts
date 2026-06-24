package streamingfile

import "sync"

var FileBufferPool = sync.Pool{
	New: func()any{
		buffer := make(chan []byte, 512)
		return &buffer
	},
}