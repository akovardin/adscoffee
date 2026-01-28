package cursor

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/theplant/relay"
)

type KeysetFinder[T any] interface {
	Find(ctx context.Context, after, before *map[string]any, orderBys []relay.OrderBy, limit int, fromEnd bool) ([]Node[T], error)
	Count(ctx context.Context) (int, error)
}

func NewKeysetAdapter[T any](finder KeysetFinder[T]) relay.ApplyCursorsFunc[T] {
	return func(ctx context.Context, req *relay.ApplyCursorsRequest) (*relay.ApplyCursorsResponse[T], error) {
		keys := lo.Map(req.OrderBys, func(item relay.OrderBy, _ int) string {
			return item.Field
		})
		if len(keys) == 0 {
			return nil, errors.New("keyset pagination requires orderBys to be set")
		}

		after, before, err := decodeKeysetCursors(req.After, req.Before, keys)
		if err != nil {
			return nil, err
		}

		skip := relay.GetSkip(ctx)

		var totalCount *int
		if !skip.TotalCount {
			count, err := finder.Count(ctx)
			if err != nil {
				return nil, err
			}
			totalCount = &count
		}

		if skip.Edges && skip.Nodes && skip.PageInfo {
			return &relay.ApplyCursorsResponse[T]{
				TotalCount: totalCount,
			}, nil
		}

		var edges []*relay.LazyEdge[T]
		if req.Limit <= 0 || (totalCount != nil && *totalCount <= 0) {
			edges = make([]*relay.LazyEdge[T], 0)
		} else {
			nodes, err := finder.Find(ctx, after, before, req.OrderBys, req.Limit, req.FromEnd)
			if err != nil {
				return nil, err
			}
			edges = make([]*relay.LazyEdge[T], len(nodes))
			for i, node := range nodes {
				edges[i] = &relay.LazyEdge[T]{
					Node: node.RelayNode(),
					Cursor: func(_ context.Context) (string, error) {
						return EncodeKeysetCursor(node, keys)
					},
				}
			}
		}

		return &relay.ApplyCursorsResponse[T]{
			LazyEdges:  edges,
			TotalCount: totalCount,
			// It would be very costly to check whether after and before really exist,
			// So it is usually not worth it. Normally, checking that it is not nil is sufficient.
			HasAfterOrPrevious: after != nil,
			HasBeforeOrNext:    before != nil,
		}, nil
	}
}

const KeysetTagKey = "relay"

// use strcut field name as key and force emit empty
var jsoniterForKeyset = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	TagKey:                 KeysetTagKey,
}.Froze()

func JSONMarshal(v any) ([]byte, error) {
	return jsoniterForKeyset.Marshal(v)
}

func JSONUnmarshal(data []byte, v any) error {
	return jsoniterForKeyset.Unmarshal(data, v)
}

func EncodeKeysetCursor(node any, keys []string) (string, error) {
	b, err := JSONMarshal(node)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal node to JSON")
	}

	m := make(map[string]any)
	if err := JSONUnmarshal(b, &m); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal node JSON to map")
	}

	keysMap := lo.SliceToMap(keys, func(key string) (string, bool) {
		return key, true
	})
	for k := range keysMap {
		if _, ok := m[k]; !ok {
			return "", errors.Errorf("required key %q not found in node when encoding cursor", k)
		}
	}
	for k := range m {
		if _, ok := keysMap[k]; !ok {
			delete(m, k)
		}
	}

	b, err = JSONMarshal(m)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal filtered cursor map to JSON")
	}
	return string(b), nil
}

func DecodeKeysetCursor(cursor string, keys []string) (map[string]any, error) {
	var m map[string]any
	if err := JSONUnmarshal([]byte(cursor), &m); err != nil {
		return nil, errors.Wrap(err, "failed to parse cursor JSON")
	}

	if len(m) != len(keys) {
		return nil, errors.Errorf("invalid cursor: has %d keys, but %d keys are expected", len(m), len(keys))
	}

	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return nil, errors.Errorf("required key %q not found in cursor", k)
		}
	}
	return m, nil
}

func decodeKeysetCursors(after, before *string, keys []string) (afterKeyset, beforeKeyset *map[string]any, err error) {
	if after != nil && before != nil && *after == *before {
		return nil, nil, errors.New("invalid pagination: after and before cursors are identical")
	}
	if after != nil {
		m, err := DecodeKeysetCursor(*after, keys)
		if err != nil {
			return nil, nil, err
		}
		afterKeyset = &m
	}
	if before != nil {
		m, err := DecodeKeysetCursor(*before, keys)
		if err != nil {
			return nil, nil, err
		}
		beforeKeyset = &m
	}
	return afterKeyset, beforeKeyset, nil
}
