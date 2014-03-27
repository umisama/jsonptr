package jsonptr

import (
	"bytes"
	"testing"
)

var JSONexample = []byte(`{
		"foo": ["bar", "baz"],
		"": 0,
		"a/b": 1,
		"c%d": 2,
		"e^f": 3,
		"g|h": 4,
		"i\\j": 5,
		"k\"l": 6,
		" ": 7,
		"m~n": 8
}`)

var TestCasesForNormal = map[string]string{
	`/foo`:   `["bar","baz"]`,
	`/foo/0`: `"bar"`,
	`/`:      `0`,
	`/a~1b`:  `1`,
	`/c%d`:   `2`,
	`/e^f`:   `3`,
	`/g|h`:   `4`,
	`/i\\j`:  `5`,
	`/k\"l`:  `6`,
	`/ `:     `7`,
	`/m~0n`:  `8`,
}

var TestCasesForURIEncoded = map[string]string{
	`#/foo`:   `["bar","baz"]`,
	`#/foo/0`: `"bar"`,
	`#/`:      `0`,
	`#/a~1b`:  `1`,
	`#/c%25d`: `2`,
	`#/e%5Ef`: `3`,
	`#/g%7Ch`: `4`,
	`#/i%5Cj`: `5`,
	`#/k%22l`: `6`,
	`#/%20`:   `7`,
	`#/m~0n`:  `8`,
}

func TestNormal(t *testing.T) {
	for k, v := range TestCasesForNormal {
		res, err := Find(JSONexample, k)
		if err != nil {
			t.Error("fail on", k, "with", err)
			return
		}

		if !bytes.Equal(res, []byte(v)) {
			t.Error("fail on", k, string(res))
		}
	}
	return
}

func TestURIEncoded(t *testing.T) {
	for k, v := range TestCasesForURIEncoded {
		res, err := Find(JSONexample, k)
		if err != nil {
			t.Error("fail on", k, "with", err)
			return
		}

		if !bytes.Equal(res, []byte(v)) {
			t.Error("fail on", k, string(res))
		}
	}
	return
}
