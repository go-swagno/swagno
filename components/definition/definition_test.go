package definition_test

import (
	"database/sql"
	"log"
	"testing"
)

func Test_F(t *testing.T) {
	i := sql.NullFloat64{
		Float64: 0,
		Valid:   true,
	}
	v, err := i.Value()
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	log.Println(v)
}
