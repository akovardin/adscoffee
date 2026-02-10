package sessions

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Session struct {
	Value string

	expiry time.Time
}

func (s Session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

type Sessions struct {
	sessions sync.Map
}

func New() *Sessions {
	return &Sessions{}
}

// TODO: в ключе нужно использовать слот
func (s *Sessions) LoadWithExpire(r *http.Request) (Session, bool) {
	session, ok := s.LoadWithoutExpire(r)

	if session.isExpired() {
		return Session{}, false
	}

	return session, ok
}

func (s *Sessions) LoadWithoutExpire(r *http.Request) (Session, bool) {
	token := s.identifier(r)
	raw, ok := s.sessions.Load(token)
	if !ok {
		return Session{}, false
	}

	session, ok := raw.(Session)

	return session, ok
}

func (s *Sessions) Start(r *http.Request, value string) error {
	token := s.identifier(r)
	expires := time.Now().Add(10 * time.Minute)

	s.sessions.Store(token, Session{
		Value:  value,
		expiry: expires,
	})

	return nil
}

func (s *Sessions) identifier(r *http.Request) string {
	agent := r.UserAgent()
	ip := forwarded(address(r), r)
	data := agent + ip
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

func address(r *http.Request) string {
	addr := r.RemoteAddr
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}

	return ip
}

func forwarded(ip string, r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0])
		}
	}

	return ip
}
