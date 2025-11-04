package log

import (
	"iter"
	"log/slog"
)

func mergeGroupAttrChildren(groupKey string, children ...slog.Attr) attrMerger {
	return func(prevAttrs iter.Seq[slog.Attr]) iter.Seq[slog.Attr] {
		return func(yield func(slog.Attr) bool) {
			var found bool
			for a := range prevAttrs {
				if a.Key != groupKey {
					if !yield(a) {
						return
					}
					continue
				}

				found = true

				if a.Value.Kind() != slog.KindGroup {
					// wanted parent is found, but it is a scalar attribute
					if !yield(a) {
						return
					}
					continue
				}

				existent := a.Value.Group()
				existent = append(existent, children...)
				if !yield(slog.GroupAttrs(groupKey, existent...)) {
					return
				}
			}

			if !found {
				_ = yield(slog.GroupAttrs(groupKey, children...))
			}
		}
	}
}
