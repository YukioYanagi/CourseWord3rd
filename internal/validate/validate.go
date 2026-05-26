package validate

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

var ErrInvalidPayload = errors.New("payload does not match declared format")

func JSON(data []byte) error {
	if !json.Valid(data) {
		return ErrInvalidPayload
	}
	return nil
}

func XML(data []byte) error {
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		_, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return ErrInvalidPayload
		}
	}
}

func SOAP(data []byte) error {
	if err := XML(data); err != nil {
		return err
	}
	s := strings.ToLower(string(data))
	if !strings.Contains(s, "envelope") {
		return ErrInvalidPayload
	}
	return nil
}
