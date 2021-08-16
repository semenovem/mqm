package mqpro

import (
  "bytes"
  "encoding/xml"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "strings"
)

// Rfh2Unmarshal получение заголовков rfh2
func (c *Mqconn) Rfh2Unmarshal(b []byte) ([]*MQRFH2, error) {
  var (
    tot []*MQRFH2
    rfh *MQRFH2
    err error
    ofs int32
  )

  for ofs < int32(len(b)) {
    rfh, err = rfh2ParseHeader(b[ofs:])
    if err != nil {
      return nil, err
    }
    if rfh == nil {
      break
    }
    tot = append(tot, rfh)
    ofs += rfh.StrucLength
  }

  return tot, nil
}

func rfh2ParseHeader(b []byte) (*MQRFH2, error) {
  if len(b) < 4 {
    return nil, nil
  }

  if !bytes.Equal([]byte(StructId), b[:4]) {
    return nil, nil
  }

  if int32(len(b)) < ibmmq.MQRFH_STRUC_LENGTH_FIXED_2 {
    return nil, ErrFormatRFH2
  }

  h := &MQRFH2{}
  var err error

  h.StructId = string(b[:4])
  h.Version = int32(endian.Uint32(b[4:8]))
  h.StrucLength = int32(endian.Uint32(b[8:12]))
  h.Encoding = int32(endian.Uint32(b[12:16]))
  h.CodedCharSetId = int32(endian.Uint32(b[16:20]))
  h.Format = strings.TrimRight(string(b[20:28]), " ")
  h.Flags = int32(endian.Uint32(b[28:32]))
  h.NameValueCCSID = int32(endian.Uint32(b[32:36]))

  if int32(len(b)) < h.StrucLength {
    return nil, ErrFormatRFH2
  }
  err = rfh2ParseData(b[36:h.StrucLength], h)
  if err != nil {
    return nil, err
  }

  return h, nil
}

// Обработка пар NameValueLength NameValueData
// https://www.ibm.com/docs/en/ibm-mq/9.0?topic=mqrfh2-namevaluelength-mqlong
func rfh2ParseData(buf []byte, rfh *MQRFH2) error {
  ofs := 0

  for ofs+4 < len(buf) {
    l := int(endian.Uint32(buf[ofs : ofs+4]))
    ofs += 4

    if len(buf) < l+ofs {
      return ErrParseRfh2
    }

    b := buf[ofs : l+ofs]
    ofs += l
    m, err := rfh2ParseXml(b)
    if err != nil {
      return err
    }

    rfh.RawXml = append(rfh.RawXml, bytes.TrimRight(b, "\x00"))
    rfh.NameValues = append(rfh.NameValues, m)
  }

  return nil
}

func rfh2ParseXml(buf []byte) (map[string]interface{}, error) {
  m := rfh2Xml{}
  err := xml.Unmarshal(buf, &m)
  if err != nil {
    return nil, err
  }

  mm := make(map[string]interface{})
  for n, v := range m.m {
    mm[n] = v
  }

  return mm, nil
}

type rfh2Xml struct {
  m map[string]interface{}
}

// UnmarshalXML
// TODO сейчас без поддержки вложенности тегов
func (c *rfh2Xml) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
  c.m = make(map[string]interface{})
  key := start.Name.Local

  for {
    t, _ := d.Token()
    switch tt := t.(type) {
    case xml.StartElement:
    case xml.EndElement:
      if tt.Name == start.Name {
        return nil
      }
    case xml.CharData:
      c.m[key] = string(tt)
    }
  }
}