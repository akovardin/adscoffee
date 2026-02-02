package sessions

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessions_LoadWithExpire_ValidSession(t *testing.T) {
	// Arrange
	sessions := New()
	req := httptest.NewRequest("GET", "/", nil)
	sessionValue := "test-session-value"

	// Start a session
	err := sessions.Start(req, sessionValue)
	assert.NoError(t, err, "Starting session should not produce an error")

	// Act
	session, ok := sessions.LoadWithExpire(req)

	// Assert
	assert.True(t, ok, "Session should be found")
	assert.Equal(t, sessionValue, session.Value, "Session value should match")
	assert.False(t, session.isExpired(), "Session should not be expired")
}

func TestSessions_LoadWithExpire_ExpiredSession(t *testing.T) {
	// Arrange
	sessions := New()
	req := httptest.NewRequest("GET", "/", nil)
	sessionValue := "test-session-value"

	// Start a session with a short expiry time
	token := sessions.identifier(req)
	expires := time.Now().Add(-1 * time.Second) // Expired 1 second ago

	sessions.sessions.Store(token, Session{
		Value:  sessionValue,
		expiry: expires,
	})

	// Act
	session, ok := sessions.LoadWithExpire(req)

	// Assert
	assert.False(t, ok, "Session should not be found as it's expired")
	assert.Equal(t, Session{}, session, "Should return empty session")
}

func TestSessions_LoadWithExpire_NonExistentSession(t *testing.T) {
	// Arrange
	sessions := New()
	req := httptest.NewRequest("GET", "/", nil)

	// Act
	session, ok := sessions.LoadWithExpire(req)

	// Assert
	assert.False(t, ok, "Session should not be found")
	assert.Equal(t, Session{}, session, "Should return empty session")
}
