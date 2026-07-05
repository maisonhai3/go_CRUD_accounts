package streamingfile

import "io"

// aws is a stand-in for a real object-storage client (e.g. the AWS S3
// SDK's uploader). These sketches are about streaming a request body
// straight through to storage without buffering the whole file, so the
// client itself is stubbed — Put just drains the reader.
var aws awsClient

type awsClient struct{}

func (awsClient) Put(r io.Reader) error {
	_, err := io.Copy(io.Discard, r)
	return err
}
