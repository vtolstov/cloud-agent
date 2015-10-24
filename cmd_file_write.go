package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

var cmdFileWrite = &Command{
	Name:    "guest-file-write",
	Func:    fnFileWrite,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdFileWrite)
}

func fnFileWrite(req *Request) *Response {
	res := &Response{Id: req.Id}

	file := struct {
		Handle int    `json:"handle"`
		BufB64 string `json:"buf-b64"`
		Count  int    `json:"count,omitempty"`
	}{}

	ret := struct {
		Count int  `json:"count"`
		Eof   bool `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &file)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		if f, ok := openFiles[file.Handle]; ok {
			var buffer []byte
			buffer, err = base64.StdEncoding.DecodeString(file.BufB64)
			if err != nil {
				res.Error = &Error{Code: -1, Desc: err.Error()}
				return res
			}
			var n int
			n, err = f.Write(buffer)
			switch err {
			case nil:
				ret.Count = n
				res.Return = ret
			case io.EOF:
				ret.Count = n
				ret.Eof = true
				res.Return = ret
			default:
				res.Error = &Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
