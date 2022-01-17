package connection

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/juju/errors"
	"github.com/obgnail/http-server-toy/context"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	LineEndFlag = "\r\n"
)

type Conn struct {
	tcpConn *net.TCPConn
	reader  *bufio.Reader
	once    sync.Once
}

func NewConn(tcpConn *net.TCPConn) *Conn {
	return &Conn{tcpConn: tcpConn, reader: bufio.NewReader(tcpConn)}
}

func (c *Conn) String() string {
	return fmt.Sprintf("%s <===> %s", c.GetRemoteAddr(), c.GetLocalAddr())
}

func (c *Conn) Close() {
	c.once.Do(c.close)
}

func (c *Conn) close() {
	if c.tcpConn != nil {
		c.tcpConn.Close()
	}
}

func (c *Conn) GetRemoteAddr() (addr string) {
	return c.tcpConn.RemoteAddr().String()
}

func (c *Conn) GetLocalAddr() (addr string) {
	return c.tcpConn.LocalAddr().String()
}

func (c *Conn) SendResponse(resp *context.Response) (err error) {
	firstLine := fmt.Sprintf("%s %d %s", resp.Proto, resp.StatusCode, resp.Status)
	buffer := bytes.NewBufferString(firstLine)
	buffer.WriteString(LineEndFlag)
	for name, hs := range resp.Header {
		headerValue := strings.Join(hs, "; ")
		buffer.WriteString(name)
		buffer.WriteString(": ")
		buffer.WriteString(headerValue)
		buffer.WriteString(LineEndFlag)
	}
	buffer.WriteString(LineEndFlag)
	buffer.Write(resp.Body)
	buffer.WriteString(LineEndFlag)
	_, err = c.tcpConn.Write(buffer.Bytes())
	if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

func (c *Conn) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *Conn) ReadLine() (line []byte, err error) {
	line, err = c.readLine()
	if err == io.EOF {
		c.Close()
	}
	return
}

func (c *Conn) readLine() (line []byte, err error) {
	for {
		l, remain, err := c.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		if line == nil && !remain {
			return l, nil
		}
		line = append(line, l...)
		if !remain {
			break
		}
	}
	return
}

func (c *Conn) GetRequest() (*context.Request, error) {
	// read first line
	firstLine, err := c.ReadLine()
	if err != nil {
		err = fmt.Errorf("read fitst line err:%s", err)
		return nil, errors.Trace(err)
	}
	method, url, proto, ok := context.ParseFirstLine(string(firstLine))
	if !ok {
		err = fmt.Errorf("parse fitst line err:%s", err)
		return nil, errors.Trace(err)
	}
	if !context.ValidMethod(method) {
		return nil, fmt.Errorf("invaild Method, %s", method)
	}
	req := &context.Request{
		Method: method,
		URL:    url,
		Proto:  proto,
	}
	req.ParsePath()

	// read headers
	header := make(map[string][]string)
	for {
		line, err := c.ReadLine()
		if err != nil {
			err = fmt.Errorf("read line err:%s", err)
			return nil, errors.Trace(err)
		}
		lineStr := string(line)

		if lineStr == "" {
			log.Debug("header fragment end")
			req.Header = header
			break
		}
		h := strings.Split(lineStr, ": ")
		if len(h) != 2 {
			err = fmt.Errorf("read header err. %s", lineStr)
			return nil, errors.Trace(err)
		}
		name, values := strings.ToLower(h[0]), strings.ToLower(h[1])
		if name == context.UserAgent {
			header[name] = []string{values}
		} else {
			header[name] = strings.Split(values, ", ")
		}
	}

	// read body
	var contentLength int
	v, ok := header[context.ContentLength]
	if !ok {
		log.Debug("request has no content-length")
		return req, nil
	}
	contentLength, err = strconv.Atoi(v[0])
	if err != nil {
		err = fmt.Errorf("strconv.Atoi content-length err, %s", v)
		return nil, err
	}
	bodyBytes := make([]byte, contentLength)
	length, err := c.Read(bodyBytes)
	if err != nil {
		if err == io.EOF {
			log.Debug("EOF")
			return nil, err
		}
	}
	if length != contentLength {
		err = fmt.Errorf("read content-length err, contentLength %d, %s", contentLength, string(bodyBytes))
		return nil, err
	}
	req.Body = bodyBytes
	return req, nil
}


