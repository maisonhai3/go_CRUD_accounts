package streamingfile

import (
	"bytes"
	"io"
	"net/http"
)

func UploadImage(w http.ResponseWriter, r *http.Request){
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	// Validate: file size
	file, fileHeader, err := r.FormFile("image")
	if err != nil || fileHeader.Size > 10 * 1024 * 1024{
		return
	}

	// Get buffer
	buffPtr := FileBufferPool.Get().(*[]byte)
	defer FileBufferPool.Put(buffPtr)

	buff := *buffPtr

	// Fill the buffer with header to determine MIME type
	headLen, err := io.ReadFull(file , buff)  // buff is holding a piece of header now
	if err != nil && err != io.ErrUnexpectedEOF{
		return
	}

	if http.DetectContentType(buff[:headLen]) != "image/jpeg"{
		return
	}

	// Re-tape the request to get a like-new
	head := bytes.NewReader(buff[:headLen])
	retapeBody := io.MultiReader(head, r.Body)

	aws.Put(retapeBody)
}