package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SSEClient struct {
	w http.ResponseWriter
	f http.Flusher
	c echo.Context
}

func NewSSEClient(c echo.Context) *SSEClient {
	return &SSEClient{
		c: c,
	}
}

func (s *SSEClient) Init() error {
	s.c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	s.c.Response().Header().Set("Cache-Control", "no-cache")
	s.c.Response().Header().Set("Connection", "keep-alive")
	s.w = s.c.Response().Writer
	flusher, ok := s.w.(http.Flusher)
	if !ok {
		return errors.New("streaming unsupported")
	}
	s.f = flusher
	return nil
}

func (s *SSEClient) Send(str string) error {
	_, err := fmt.Fprintf(s.w, "data: %s\n\n", str)
	if err != nil {
		return err
	}
	s.f.Flush()
	return nil
}

func (s *SSEClient) SendResponse(resp interface{}) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return s.Send(string(data))
}

func (s *SSEClient) Close() {
	fmt.Fprint(s.w, "data: [DONE]\n\n")
	s.f.Flush()
}
