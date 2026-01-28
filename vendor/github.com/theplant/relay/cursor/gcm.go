package cursor

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/theplant/relay"
)

func encryptGCM(gcm cipher.AEAD, plainText string) (string, error) {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "could not generate nonce")
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.RawURLEncoding.EncodeToString(cipherText), nil
}

func decryptGCM(gcm cipher.AEAD, cipherText string) (string, error) {
	decodedCipherText, err := base64.RawURLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", errors.Wrap(err, "could not decode cipher text")
	}

	nonceSize := gcm.NonceSize()
	if len(decodedCipherText) < nonceSize {
		return "", errors.New("cipher text too short")
	}

	nonce, dataCipherText := decodedCipherText[:nonceSize], decodedCipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, dataCipherText, nil)
	if err != nil {
		return "", errors.Wrap(err, "could not decrypt cipher text")
	}

	return string(plainText), nil
}

// NewGCM creates a new GCM cipher
// Concurrent safe: https://github.com/golang/go/issues/41689
func NewGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "could not create cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "could not create AEAD")
	}
	return gcm, nil
}

func GCM[T any](gcm cipher.AEAD) relay.CursorMiddleware[T] {
	return func(next relay.ApplyCursorsFunc[T]) relay.ApplyCursorsFunc[T] {
		return func(ctx context.Context, req *relay.ApplyCursorsRequest) (*relay.ApplyCursorsResponse[T], error) {
			if req.After != nil {
				decodedCursor, err := decryptGCM(gcm, *req.After)
				if err != nil {
					return nil, errors.Wrap(err, "invalid after cursor")
				}
				req.After = lo.ToPtr(decodedCursor)
			}

			if req.Before != nil {
				decodedCursor, err := decryptGCM(gcm, *req.Before)
				if err != nil {
					return nil, errors.Wrap(err, "invalid before cursor")
				}
				req.Before = lo.ToPtr(decodedCursor)
			}

			rsp, err := next(ctx, req)
			if err != nil {
				return nil, err
			}

			for _, edge := range rsp.LazyEdges {
				originalCursor := edge.Cursor
				edge.Cursor = func(ctx context.Context) (string, error) {
					cursor, err := originalCursor(ctx)
					if err != nil {
						return "", err
					}
					encryptedCursor, err := encryptGCM(gcm, cursor)
					if err != nil {
						return "", err
					}
					return encryptedCursor, nil
				}
			}

			return rsp, nil
		}
	}
}
