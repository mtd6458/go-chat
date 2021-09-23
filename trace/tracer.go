package trace

import (
	"fmt"
	"io"
)

// Tracer Tracerはコード内での出来事を記録できるオブジェクトを表すインターフェースです。
// ...Interface{}という型は、任意の型の引数を何個でも受け取ることができる。(０個でもOK
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	// 情報出力用のフィールド
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	if _, err := t.out.Write([]byte(fmt.Sprint(a...))); err != nil {
		return
	}

	if _, err := t.out.Write([]byte("\n")); err != nil {
		return
	}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}
// Off はTraceメソッドの呼び出しを無視するTracerを返します
func Off() Tracer {
	return &nilTracer{}
}
