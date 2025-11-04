package log

import (
	"iter"
	"log/slog"
	"slices"
)

type attrMerger func(iter.Seq[slog.Attr]) iter.Seq[slog.Attr]

func mergeRecordAttrs(base slog.Record, merger attrMerger) slog.Record {
	ret := slog.NewRecord(base.Time, base.Level, base.Message, base.PC)
	ret.AddAttrs(slices.Collect(merger(base.Attrs))...)
	return ret
}
