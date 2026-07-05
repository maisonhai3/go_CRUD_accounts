package streamingfile

import "sync"

var FileBufferPool = sync.Pool{
	New: func()any{
		buffer := make([]byte, 512)
		return &buffer
	},
}