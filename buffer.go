package lightdouble

import "strings"

type builderBuffer struct {
	limit   int
	builder strings.Builder
}

func newBuffer(limit int) *builderBuffer {
	ss := strings.Builder{}
	return &builderBuffer{
		limit:   limit,
		builder: ss,
	}
}

func (buffer *builderBuffer) write(data string) int {
	bs := stringToBytes(data)
	need := len(bs)
	if need+buffer.builder.Len() > int(buffer.limit) {
		return -1
	}
	buffer.builder.WriteString(data)
	return 1
}

func (buffer *builderBuffer) reset() {
	buffer.builder.Reset()
}

func (buffer *builderBuffer) getString() string {
	return buffer.builder.String()
}
