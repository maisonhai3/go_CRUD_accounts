package streamingfile

import (
	"bytes"
	"io"
	"net/http"
)

func RequestHandle(r http.Request){
	// Check if Content Type is Image
	var buffPtr = FileBufferPool.Get().(*[]byte)
	defer FileBufferPool.Put(buffPtr)

	buff := *buffPtr

	headLen, err := io.ReadFull(r.Body, buff)

	if http.DetectContentType(buff[:headLen]) != "image/jpeg" {
		return
	}

	// Re-glue
	buffReader := bytes.NewReader(buff[:headLen])
	restoreBody := io.MultiReader(buffReader, r.Body)
	
	// Forward the Body to AWS
	aws.Put(restoreBody)
	
}