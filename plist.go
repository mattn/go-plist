package plist

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

type Array []interface{}
type Dict map[string]interface{}

func Read(r io.Reader) (interface{}, error) {
	p := xml.NewDecoder(r)
	_, err := nextStart(p) // plist
	if err != nil {
		return nil, err
	}
	_, root, err := next(p) // root
	if err != nil {
		return nil, err
	}
	return root, nil
}

var xmlSpecial = map[byte]string{
	'<':  "&lt;",
	'>':  "&gt;",
	'"':  "&quot;",
	'\'': "&apos;",
	'&':  "&amp;",
}

func xmlEscape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if s, ok := xmlSpecial[c]; ok {
			b.WriteString(s)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

type valueNode struct {
	Type string `xml:"attr"`
	Body string `xml:"chardata"`
}

func next(p *xml.Decoder) (xml.Name, interface{}, error) {
	se, e := nextStart(p)
	if e != nil {
		return xml.Name{}, nil, e
	}

	var nv interface{}
	switch se.Name.Local {
	case "string":
		var s string
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		return xml.Name{}, s, nil
	case "true":
		return xml.Name{}, true, e
	case "false":
		return xml.Name{}, false, e
	case "integer":
		var s string
		var i int
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		i, e = strconv.Atoi(strings.TrimSpace(s))
		return xml.Name{}, i, e
	case "double":
		var s string
		var f float64
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		f, e = strconv.ParseFloat(strings.TrimSpace(s), 64)
		return xml.Name{}, f, e
	case "date":
		var s string
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		t, e := time.Parse("2006-01-02T15:04:05Z", s)
		return xml.Name{}, t, e
	case "data":
		var s string
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		if b, e := base64.StdEncoding.DecodeString(s); e != nil {
			return xml.Name{}, nil, e
		} else {
			return xml.Name{}, b, nil
		}
	case "value":
		_, e = nextStart(p)
		if e != nil {
			return xml.Name{}, nil, e
		}
		return next(p)
	case "name":
		_, e = nextStart(p)
		if e != nil {
			return xml.Name{}, nil, e
		}
		return next(p)
	case "dict":
		st := Dict{}

		for e == nil {
			se, e = nextStart(p)
			if e != nil {
				break
			}
			if se.Name.Local != "key" {
				return xml.Name{}, nil, errors.New("invalid key ")
			}
			var key string
			if e = p.DecodeElement(&key, &se); e != nil {
				return xml.Name{}, nil, e
			}

			// value
			var value interface{}
			_, value, e = next(p)
			if e != nil {
				break
			}
			st[key] = value
		}
		return xml.Name{}, st, nil
	case "array":
		ar := Array{}

		for {
			_, e = nextStart(p)
			if e != nil {
				break
			}
			_, value, e := next(p)
			if e != nil {
				break
			}
			ar = append(ar, value)
		}
		return xml.Name{}, ar, nil
	}

	if e = p.DecodeElement(nv, &se); e != nil {
		return xml.Name{}, nil, e
	}
	return se.Name, nv, e
}

func nextStart(p *xml.Decoder) (xml.StartElement, error) {
	for {
		t, e := p.Token()
		if e != nil {
			return xml.StartElement{}, e
		}
		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		}
	}
	panic("unreachable")
}
