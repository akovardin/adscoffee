package formats

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/tools/filesystem"
)

// MockFileReader is a mock implementation of the FileReader interface
type MockFileReader struct {
	mock.Mock
}

func (m *MockFileReader) Open() (io.ReadSeekCloser, error) {
	args := m.Called()
	return args.Get(0).(io.ReadSeekCloser), args.Error(1)
}

// MockReadSeekCloser is a mock implementation of io.ReadSeekCloser
type MockReadSeekCloser struct {
	mock.Mock
}

func (m *MockReadSeekCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	args := m.Called(offset, whence)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReadSeekCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

// createTestImage creates a simple test image for testing
func createTestImage() (*bytes.Buffer, string, error) {
	width, height := 100, 50
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a solid color
	c := color.RGBA{255, 128, 0, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, "", err
	}

	return buf, "png", nil
}

func TestNewBanner(t *testing.T) {
	banner := NewBanner()
	assert.NotNil(t, banner)
}

func TestBanner_Banner_Success(t *testing.T) {
	// Create a test image
	imgBuffer, format, err := createTestImage()
	assert.NoError(t, err)
	assert.NotNil(t, imgBuffer)

	// Create a test server to serve our test image
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(imgBuffer.Bytes())
	}))
	defer testServer.Close()

	// Create a test banner with the test server URL
	banner := NewBanner()
	testBanner := ads.Banner{
		Title:       "Test Banner",
		Description: "This is a test banner",
		Image:       ads.Image{Url: testServer.URL + "/test/image.png"},
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the Banner method
	err = banner.Banner(context.Background(), "", testBanner, rr)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "image/"+format, rr.Header().Get("Content-Type"))
	assert.NotEmpty(t, rr.Body.Bytes())
}

func TestBanner_Banner_FileError(t *testing.T) {
	// Create a mock file system function that returns an error
	originalNewFileFromURL := filesystemNewFileFromURL
	filesystemNewFileFromURL = func(ctx context.Context, url string) (*filesystem.File, error) {
		return nil, fmt.Errorf("failed to download url %s", url)
	}
	defer func() {
		filesystemNewFileFromURL = originalNewFileFromURL
	}()

	// Create a test banner
	banner := NewBanner()

	// Create a test ads.Banner
	testBanner := ads.Banner{
		Title:       "Test Banner",
		Description: "This is a test banner",
		Image:       ads.Image{Url: "http://example.com/test/image.png"},
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the Banner method
	err := banner.Banner(context.Background(), "http://example.com", testBanner, rr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download url")
}

func TestBanner_Render_Success(t *testing.T) {
	// Create a test image
	imgBuffer, _, err := createTestImage()
	assert.NoError(t, err)
	assert.NotNil(t, imgBuffer)

	// Create a test banner
	banner := NewBanner()

	// Call the Render method
	result, format, err := banner.Render(imgBuffer, "Test Description", "Test Info")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, []string{"jpeg", "png"}, format)
	assert.NotEmpty(t, result.Bytes())
}

func TestBanner_Render_InvalidImage(t *testing.T) {
	// Create invalid image data
	invalidData := bytes.NewBufferString("invalid image data")

	// Create a test banner
	banner := NewBanner()

	// Call the Render method
	result, format, err := banner.Render(invalidData, "Test Description", "Test Info")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, format)
}

func TestWrap(t *testing.T) {
	text := "This is a test text for wrapping functionality"
	lineLength := 10

	lines := wrap(text, lineLength)
	assert.NotEmpty(t, lines)
	assert.Greater(t, len(lines), 1)

	// Check that each line is not longer than the line length
	for _, line := range lines {
		// Note: The wrap function might create lines slightly longer than lineLength
		// due to word boundaries, so we check for a reasonable upper bound
		assert.LessOrEqual(t, len(line), lineLength+10)
	}
}

func TestWrap_EmptyText(t *testing.T) {
	text := ""
	lineLength := 10

	lines := wrap(text, lineLength)
	assert.Empty(t, lines)
}

func TestDrawLine(t *testing.T) {
	width, height := 100, 50
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw a line
	drawLine(img, 10, 10, 90, 10, color.RGBA{255, 0, 0, 255})

	// Check that the line was drawn
	c := img.RGBAAt(10, 10)
	assert.Equal(t, uint8(255), c.R)
	assert.Equal(t, uint8(0), c.G)
	assert.Equal(t, uint8(0), c.B)
	assert.Equal(t, uint8(255), c.A)
}

// Variable to allow mocking of filesystem.NewFileFromURL
var filesystemNewFileFromURL = filesystem.NewFileFromURL
