package plist

import (
	"encoding/base64"
	"encoding/xml"
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
		if err != io.EOF {
			return nil, err
		}
	}
	return root, nil
}

func next(p *xml.Decoder) (xml.Name, interface{}, error) {
	sei, e := nextStart(p)
	if e != nil {
		return xml.Name{}, nil, e
	}
	if _, ok := sei.(xml.EndElement); ok {
		return xml.Name{}, nil, nil
	}
	se := sei.(xml.StartElement)

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
		var i int64
		if e = p.DecodeElement(&s, &se); e != nil {
			return xml.Name{}, nil, e
		}
		i, e = strconv.ParseInt(strings.TrimSpace(s), 10, 64)
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
		s = strings.TrimSpace(s)
		s = strings.Replace(s, "\t", "", -1)
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		if b, e := base64.StdEncoding.DecodeString(s); e != nil {
			return xml.Name{}, nil, e
		} else {
			return xml.Name{}, b, nil
		}
	case "dict":
		st := Dict{}
		n := 0
		for e == nil {
			// key
			sei, e = nextStart(p)
			if e != nil {
				break
			}
			if ee, ok := sei.(xml.EndElement); ok {
				switch ee.Name.Local {
				case "true":
					continue
				case "false":
					continue
				}
				break
			}
			se := sei.(xml.StartElement)
			var key string
			if e = p.DecodeElement(&key, &se); e != nil {
				return xml.Name{}, nil, e
			}

			// value
			var value interface{}
			_, value, e = next(p)
			st[key] = value
			n++
		}
		return xml.Name{}, st, e
	case "array":
		ar := Array{}
		for {
			_, value, e := next(p)
			if e != nil {
				return xml.Name{}, ar, e
			}
			ar = append(ar, value)
		}
		return xml.Name{}, ar, nil
	}

	var nv interface{}
	if e = p.DecodeElement(&nv, &se); e != nil {
		return xml.Name{}, nil, e
	}
	return se.Name, nv, e
}

func nextStart(p *xml.Decoder) (interface{}, error) {
	for {
		t, e := p.Token()
		if e != nil {
			return xml.StartElement{}, e
		}
		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		case xml.EndElement:
			return t, nil
		}
	}
	panic("unreachable")
}
