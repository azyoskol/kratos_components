package gin

import (
	"context"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	var (
		srv = NewServer(
			gin.New(),
			Addr(":8080"),
		)
		ch = make(chan string, 1)
		wg sync.WaitGroup
	)
	wg.Add(1)

	srv.GET("/ping", func(c *gin.Context) {
		defer wg.Done()
		ch <- "pong"
		c.String(200, "pong")
	})

	go func() {
		srv.Start(context.Background()) // nolint: errcheck
	}()

	resp, err := http.Get("http://localhost:8080/ping")

	wg.Wait()

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close() // nolint: errcheck
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "pong", string(body))

	assert.Equal(t, "pong", <-ch)

	err = srv.Stop(context.Background())
	assert.NoError(t, err)
}
