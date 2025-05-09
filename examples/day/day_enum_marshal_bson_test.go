package day_test

import (
	"fmt"
	"testing"

	"github.com/klippa-app/go-enum/examples/day"
	"go.mongodb.org/mongo-driver/bson"
)

type testData struct {
	input []byte
	want  struct {
		err    error
		output marshallableStruct
	}
}

type marshallableStruct struct {
	Dag day.Day
}

func TestDagBSON(t *testing.T) {

	tests := []testData{
		{
			input: []byte{21, 0, 0, 0, 2, 100, 97, 103, 0, 7, 0, 0, 0, 70, 82, 73, 68, 65, 89, 0, 0}, // {dag: FRIDAY}
			want: struct {
				err    error
				output marshallableStruct
			}{
				err: nil,
				output: marshallableStruct{
					Dag: day.Friday,
				},
			},
		},
		{
			input: []byte{21, 0, 0, 0, 2, 100, 97, 103, 0, 7, 0, 0, 0, 77, 79, 78, 68, 65, 89, 0, 0}, // {dag: MONDAY}
			want: struct {
				err    error
				output marshallableStruct
			}{
				err: nil,
				output: marshallableStruct{
					Dag: day.Monday,
				},
			},
		},
		{
			input: []byte{21, 0, 0, 0, 2, 100, 97, 103, 0, 7, 0, 0, 0, 109, 111, 110, 100, 97, 121, 0, 0}, // {dag: monday}
			want: struct {
				err    error
				output marshallableStruct
			}{
				err: fmt.Errorf("error decoding key dag: monday is not a valid Day"),
				output: marshallableStruct{
					Dag: day.Unknown,
				},
			},
		},
	}

	for i := range tests {

		test := tests[i]

		data := test.input

		var res marshallableStruct = marshallableStruct{}

		err := bson.Unmarshal(data, &res)
		if test.want.err == nil {
			if err != nil {
				t.Error("expected no error got:", err)
			}
		} else if err != nil {
			if err.Error() != test.want.err.Error() {
				t.Error("expected", test.want.err, "got", err)
				continue
			}
			continue
		} else {
			t.Error("expected", test.want.err, "got nil")
		}

		if res.Dag != test.want.output.Dag {
			t.Error("invalid day", res.Dag, "expected:", test.want.output.Dag)
			continue
		}
	}

}
