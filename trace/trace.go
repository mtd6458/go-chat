package trace

// Tracer Tracerはコード内での出来事を記録できるオブジェクトを表すインターフェースです。
// ...Interface{}という型は、任意の型の引数を何個でも受け取ることができる。(０個でもOK
type Tracer interface {
	Trace(...interface{})
}
