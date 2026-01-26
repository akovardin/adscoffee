package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockManager struct {
	mountCalled bool
}

func (m *mockManager) Mount(router *chi.Mux) {
	m.mountCalled = true
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func TestNew(t *testing.T) {
	mockMgr := &mockManager{}

	server := &Server{
		srv:     &http.Server{Addr: addr},
		manager: mockMgr,
	}

	assert.NotNil(t, server)
	assert.Equal(t, mockMgr, server.manager)
	assert.NotNil(t, server.srv)
	assert.Equal(t, addr, server.srv.Addr)
}

func TestServer_Start_Success(t *testing.T) {
	mockMgr := &mockManager{}

	server := &Server{
		srv:     &http.Server{Addr: addr},
		manager: mockMgr,
	}

	ctx := context.Background()

	err := server.Start(ctx)

	assert.NoError(t, err)
	assert.True(t, mockMgr.mountCalled)

	time.Sleep(100 * time.Millisecond)

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	assert.NoError(t, err)
}

func TestServer_Shutdown(t *testing.T) {
	oldServeMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()
	defer func() {
		http.DefaultServeMux = oldServeMux
	}()

	mockMgr := &mockManager{}

	server := &Server{
		srv:     &http.Server{Addr: addr},
		manager: mockMgr,
	}

	ctx := context.Background()
	err := server.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	assert.NoError(t, err)
}

func TestServer_Integration(t *testing.T) {
	oldServeMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()
	defer func() {
		http.DefaultServeMux = oldServeMux
	}()

	mockMgr := &mockManager{}

	server := &Server{
		srv:     &http.Server{Addr: addr},
		manager: mockMgr,
	}

	ctx := context.Background()
	err := server.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost" + addr)
	if err == nil {
		resp.Body.Close()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	assert.NoError(t, err)

	assert.True(t, mockMgr.mountCalled)
}
