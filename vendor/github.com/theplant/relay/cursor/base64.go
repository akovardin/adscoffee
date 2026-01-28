package cursor

import (
	"context"
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/theplant/relay"
)

func Base64[T any](next relay.ApplyCursorsFunc[T]) relay.ApplyCursorsFunc[T] {
	return func(ctx context.Context, req *relay.ApplyCursorsRequest) (*relay.ApplyCursorsResponse[T], error) {
		if req.After != nil {
			cursor, err := base64.RawURLEncoding.DecodeString(*req.After)
			if err != nil {
				return nil, errors.Wrap(err, "invalid after cursor")
			}
			req.After = lo.ToPtr(string(cursor))
		}

		if req.Before != nil {
			cursor, err := base64.RawURLEncoding.DecodeString(*req.Before)
			if err != nil {
				return nil, errors.Wrap(err, "invalid before cursor")
			}
			req.Before = lo.ToPtr(string(cursor))
		}

		rsp, err := next(ctx, req)
		if err != nil {
			return nil, err
		}

		// Encrypt the cursor
		for _, edge := range rsp.LazyEdges {
			originalCursor := edge.Cursor
			edge.Cursor = func(ctx context.Context) (string, error) {
				cursor, err := originalCursor(ctx)
				if err != nil {
					return "", err
				}
				return base64.RawURLEncoding.EncodeToString([]byte(cursor)), nil
			}
		}

		return rsp, nil
	}
}
