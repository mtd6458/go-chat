package trace

import (
	"bytes"
	"testing"
)

// t *testing.T型の引数を一つ受け取る関数は全てユニットテストとみなされる。
func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは")
		if buf.String() != "こんにちは\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
}

func TestOff(t *testing.T)  {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}
