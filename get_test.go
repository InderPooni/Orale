package orale_test

import (
	"testing"

	orale "github.com/RobertWHurst/Orale"
)

func newTestLoaderSingleValues() *orale.Loader {
	return &orale.Loader{
		FlagValues: map[string][]any{
			"a": {"1"},
			"b": {"2"},
			"c": {"3"},
			"d": {"4"},
		},
		EnvironmentValues: map[string][]any{
			"b": {"5"},
			"e": {"6"},
		},
		ConfigurationFiles: []*orale.File{
			{
				Path: "path/to/other/file-2.toml",
				Values: map[string][]any{
					"d": {"9"},
					"g": {"10"},
				},
			},
			{
				Path: "path/to/file-1.toml",
				Values: map[string][]any{
					"c": {"7"},
					"f": {"8"},
					"g": {"9"},
					"h": {"10"},
				},
			},
		},
	}
}

func newTestLoaderMultiValues() *orale.Loader {
	return &orale.Loader{
		FlagValues: map[string][]any{
			"a": {"1", "2"},
			"b": {"3", "4"},
			"c": {"5", "6"},
			"d": {"7", "8"},
		},
		EnvironmentValues: map[string][]any{
			"b": {"9", "10"},
			"e": {"11", "12"},
		},
		ConfigurationFiles: []*orale.File{},
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("should correctly resolve values into struct", func(t *testing.T) {
		t.Parallel()

		type TestStruct struct {
			A string `config:"a"`
			B string `config:"b"`
			C string `config:"c"`
			D string `config:"d"`
			F string `config:"f"`
			G string `config:"g"`
			H string `config:"h"`
			I string `config:"i"`
		}

		testStruct := TestStruct{}

		conf := newTestLoaderSingleValues()
		if err := conf.Get("", &testStruct); err != nil {
			t.Fatal(err)
		}

		if testStruct.A != "1" {
			t.Fatalf("expected A to be 1, got %s", testStruct.A)
		}
		if testStruct.B != "2" {
			t.Fatalf("expected B to be 2, got %s", testStruct.B)
		}
		if testStruct.C != "3" {
			t.Fatalf("expected C to be 3, got %s", testStruct.C)
		}
		if testStruct.D != "4" {
			t.Fatalf("expected D to be 4, got %s", testStruct.D)
		}
		if testStruct.F != "8" {
			t.Fatalf("expected F to be 8, got %s", testStruct.F)
		}
		if testStruct.G != "10" {
			t.Fatalf("expected G to be 10, got %s", testStruct.G)
		}
		if testStruct.H != "10" {
			t.Fatalf("expected H to be 10, got %s", testStruct.H)
		}
		if testStruct.I != "" {
			t.Fatalf("expected I to be empty, got %s", testStruct.I)
		}
	})

	t.Run("should correctly resolve multi values into struct", func(t *testing.T) {
		type TestStruct struct {
			A []string `config:"a"`
			B []string `config:"b"`
			C []string `config:"c"`
			D []string `config:"d"`
			E []string `config:"e"`
		}

		testStruct := TestStruct{}

		conf := newTestLoaderMultiValues()
		if err := conf.Get("", &testStruct); err != nil {
			t.Fatal(err)
		}

		if len(testStruct.A) != 2 {
			t.Fatalf("expected A to have 2 values, got %d", len(testStruct.A))
		}
		if testStruct.A[0] != "1" {
			t.Fatalf("expected A[0] to be 1, got %s", testStruct.A[0])
		}
		if testStruct.A[1] != "2" {
			t.Fatalf("expected A[1] to be 2, got %s", testStruct.A[1])
		}
		if len(testStruct.B) != 2 {
			t.Fatalf("expected B to have 2 values, got %d", len(testStruct.B))
		}
		if testStruct.B[0] != "3" {
			t.Fatalf("expected B[0] to be 3, got %s", testStruct.B[0])
		}
		if testStruct.B[1] != "4" {
			t.Fatalf("expected B[1] to be 4, got %s", testStruct.B[1])
		}
		if len(testStruct.C) != 2 {
			t.Fatalf("expected C to have 2 values, got %d", len(testStruct.C))
		}
		if testStruct.C[0] != "5" {
			t.Fatalf("expected C[0] to be 5, got %s", testStruct.C[0])
		}
		if testStruct.C[1] != "6" {
			t.Fatalf("expected C[1] to be 6, got %s", testStruct.C[1])
		}
		if len(testStruct.D) != 2 {
			t.Fatalf("expected D to have 2 values, got %d", len(testStruct.D))
		}
		if testStruct.D[0] != "7" {
			t.Fatalf("expected D[0] to be 7, got %s", testStruct.D[0])
		}
		if testStruct.D[1] != "8" {
			t.Fatalf("expected D[1] to be 8, got %s", testStruct.D[1])
		}
		if len(testStruct.E) != 2 {
			t.Fatalf("expected E to have 2 values, got %d", len(testStruct.E))
		}
		if testStruct.E[0] != "11" {
			t.Fatalf("expected E[0] to be 11, got %s", testStruct.E[0])
		}
		if testStruct.E[1] != "12" {
			t.Fatalf("expected E[1] to be 12, got %s", testStruct.E[1])
		}
	})
}
