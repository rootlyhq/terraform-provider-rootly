package client

import (
	"bytes"
	"io"

	"github.com/google/jsonapi"
)

func MarshalData(entity interface{}) (io.Reader, error) {
	buffer := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buffer, entity); err != nil {
		return nil, err
	}
	return buffer, nil
}

func UnmarshalData(data io.ReadCloser, entity interface{}) (interface{}, error) {
	if err := jsonapi.UnmarshalPayload(data, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
