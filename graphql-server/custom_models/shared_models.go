package custom_models

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

const rfc3339fulltime = "2006-02-01 15:04:05"

func MarshalID(id int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
	})
}

func UnmarshalID(v interface{}) (int, error) {
	id, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids must be strings")
	}
	i, e := strconv.Atoi(id)
	return int(i), e
}

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(t.Format(rfc3339fulltime)))
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	tmpStr, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("%T is in the wrong format", v)
	}

	time, err := time.Parse(rfc3339fulltime, tmpStr)
	if err != nil {
		return time, err
	}

	return time, nil
}
