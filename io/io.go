package io

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/logging"
	"reflect"
)

type DataBytesBuffer struct {
	*bytes.Buffer
}

func NewDataBuffer(buf []byte) *DataBytesBuffer {
	data := &DataBytesBuffer{}
	data.Buffer = bytes.NewBuffer(buf)
	return data
}

func (data *DataBytesBuffer) ReadUTF() (string, error) {
	var utfLen uint16
	err := binary.Read(data, binary.BigEndian, &utfLen)
	if err != nil {
		return "", err
	}
	bt := make([]byte, utfLen)
	realLen, err := data.Read(bt)
	if err != nil {
		return "", err
	}
	if uint16(realLen) != utfLen {
		err = errors.New("damaged package")
	}

	return string(bt), err
}

func (data *DataBytesBuffer) WriteUTF(msg string) error {
	utfLen, err := strLen(msg)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.BigEndian, utfLen)
	if err != nil {
		return err
	}
	err = binary.Write(data, binary.BigEndian, []byte(msg))
	return err
}

// ReadData Like binary.Read, but also read string.
// The field of struct can have ignored tag, it accepts the name of another field, if this value is false, then this field is not read.
// It also works on DataBytesBuffer.WriteData
func (data *DataBytesBuffer) ReadData(d ...interface{}) (err error) {
	for _, e := range d {
		v := reflect.ValueOf(e)
		k := v.Kind()
		if k != reflect.Ptr {
			return errors.New("invalid type. must be ptr")
		}
		v = v.Elem()
		k = v.Kind()
		if k == reflect.String {
			str, err := data.ReadUTF()
			if err != nil {
				return err
			}
			*e.(*string) = str
		} else if k == reflect.Struct {
			t := v.Type()
			l := v.NumField()
			var ignore []string
			for i := 0; i < l; i++ {
				v, t := v.Field(i), t.Field(i)
				if ignoreBool := has(ignore, t.Name); v.CanSet() && v.CanAddr() && !ignoreBool {
					if v.Kind() == reflect.String {
						str, err := data.ReadUTF()
						if err != nil {
							return err
						}
						v.SetString(str)
					} else {
						err := binary.Read(data, binary.BigEndian, v.Addr().Interface())
						if name, ok := t.Tag.Lookup("ignore"); v.Kind() == reflect.Bool && ok {
							if !v.Interface().(bool) {
								ignore = append(ignore, name)
							}
						}
						if err != nil {
							return err
						}
					}
				} else {
					if ignoreBool {
						continue
					}
					if v.Kind() == reflect.String {
						_, err := data.ReadUTF()
						if err != nil {
							return err
						}
					} else {
						size := v.Type().Size()
						data.Next(int(size))
					}
				}
			}
		} else {
			err = binary.Read(data, binary.BigEndian, e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteData Multi-function write data.
// For example:
//     data := io.NewDataBuffer([]byte{})
//     data.WriteData(
//         0, // int will change into int32
//         byte(123),
//         "str" // this will call DataBytesBuffer.WriteUTF
//    )
//    data.WriteData(1)
//    data.WriteData("str") // You can simply write basic data. Of course, you can also write a string
//    // You can also write into the structure. Note that there can be no int or uint and nested struct fields.
//    // No unexported field
//    data.WriteData(&struct{
//        SomeIntData int32
//        SomeStrData string
//    }{
//        1,
//        "str",
//    })
//
func (data *DataBytesBuffer) WriteData(d ...interface{}) error {
	for _, d := range d {
		v := reflect.ValueOf(d)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.String:
			err := data.WriteUTF(v.String())
			if err != nil {
				return err
			}
		case reflect.Int:
			err := binary.Write(data, binary.BigEndian, int32(v.Int()))
			if err != nil {
				return err
			}
		case reflect.Uint:
			err := binary.Write(data, binary.BigEndian, uint32(v.Uint()))
			if err != nil {
				return err
			}
		case reflect.Struct:
			t := v.Type()
			l := v.NumField()
			var ignore []string
			for i := 0; i < l; i++ {
				v, t := v.Field(i), t.Field(i)
				if ignoreBool := has(ignore, t.Name); v.CanSet() && !ignoreBool {
					if v.Kind() == reflect.String {
						err := data.WriteUTF(v.String())
						if err != nil {
							return err
						}
					} else {
						err := binary.Write(data, binary.BigEndian, v.Interface())
						if name, ok := t.Tag.Lookup("ignore"); v.Kind() == reflect.Bool && ok {
							if b := !v.Interface().(bool); b {
								err = binary.Write(data, binary.BigEndian, b)
								ignore = append(ignore, name)
							}
						}
						if err != nil {
							return err
						}
					}
				}
			}
		default:
			err := binary.Write(data, binary.BigEndian, d)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteGzipData The data is compressed by gzip and written into data.
// str is the header of this data. d will write data through WriteData.
func (data *DataBytesBuffer) WriteGzipData(str string, d ...interface{}) error {
	buf := NewDataBuffer([]byte{})
	err := buf.WriteData(d...)
	if err != nil {
		return err
	}
	dataBuf := buf.Bytes()
	buf.Reset()
	g := gzip.NewWriter(buf)
	_, err = g.Write(dataBuf)
	if err != nil {
		return err
	}
	err = g.Flush()
	if err != nil {
		return err
	}
	err = g.Close()
	if err != nil {
		return err
	}
	err = data.WriteUTF(str)
	if err != nil {
		return err
	}
	dataBuf = buf.Bytes()
	err = data.WriteData(len(dataBuf))
	if err != nil {
		return err
	}
	_, err = data.Write(dataBuf)
	//fmt.Printf("hex: %v\n", dataBuf)
	return err
}

func (data *DataBytesBuffer) ReadGzipData() (header string, d []byte, err error) {
	g, _ := gzip.NewReader(data)
	var length int32
	d = make([]byte, length)
	_, err = g.Read(d)
	if err != nil {
		logging.LogErr(err)
		return "", nil, err
	}
	return
}

func strLen(str string) (uint16, error) {
	var strLen uint16
	if strLen = uint16(len(str)); strLen > 65535 {
		return 0, errors.New("the str len is too big")
	}
	return strLen, nil
}

func has(s []string, str string) bool {
	for _, find := range s {
		if find == str {
			return true
		}
	}

	return false
}
