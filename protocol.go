package uart

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
)

type (
	Protocol interface {
		NewEncoder(io.Writer) Encoder
		NewDecoder(io.Reader) Decoder
	}

	Encoder interface {
		Encode(interface{}) error
	}

	Decoder interface {
		Decode(interface{}) error
	}

	gobProtocol  struct{}
	xmlProtocol  struct{}
	jsonProtocol struct{}
)

var (
	GobProtocol  = &gobProtocol{}
	XmlProtocol  = &xmlProtocol{}
	JsonProtocol = &jsonProtocol{}
)

func (this *gobProtocol) NewEncoder(dst io.Writer) Encoder {
	return gob.NewEncoder(dst)
}

func (this *gobProtocol) NewDecoder(src io.Reader) Decoder {
	return gob.NewDecoder(src)
}

func (this *xmlProtocol) NewEncoder(dst io.Writer) Encoder {
	return xml.NewEncoder(dst)
}

func (this *xmlProtocol) NewDecoder(src io.Reader) Decoder {
	return xml.NewDecoder(src)
}

func (this *jsonProtocol) NewEncoder(dst io.Writer) Encoder {
	return json.NewEncoder(dst)
}

func (this *jsonProtocol) NewDecoder(src io.Reader) Decoder {
	return json.NewDecoder(src)
}
