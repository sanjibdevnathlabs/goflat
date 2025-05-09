package goflat

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		given   string
		options *Options
		want    map[string]interface{}
	}{
		// test with different primitives and upper/lower case
		// String: 'world',
		// Number: 1234.99,
		// Boolean: true,
		// null: null,
		{
			`{"hello": "world"}`,
			nil,
			map[string]interface{}{"hello": "world"},
		},
		{
			`{"Hello": "world"}`,
			nil,
			map[string]interface{}{"Hello": "world"},
		},
		{
			`{"hello": 1234.99}`,
			nil,
			map[string]interface{}{"hello": 1234.99},
		},
		{
			`{"hello": true}`,
			nil,
			map[string]interface{}{"hello": true},
		},
		{
			`{"hello": null}`,
			nil,
			map[string]interface{}{"hello": nil},
		},
		// nested once
		{
			`{"hello":{}}`,
			nil,
			map[string]interface{}{"hello": map[string]interface{}{}},
		},
		{
			`{"hello":{"world":"good morning"}}`,
			nil,
			map[string]interface{}{"hello.world": "good morning"},
		},
		{
			`{"Hello":{"world":"good morning"}}`,
			nil,
			map[string]interface{}{"Hello.world": "good morning"},
		},
		{
			`{"hello":{"World":"good morning"}}`,
			nil,
			map[string]interface{}{"hello.World": "good morning"},
		},
		{
			`{"Hello":{"World":"good morning"}}`,
			nil,
			map[string]interface{}{"Hello.World": "good morning"},
		},
		{
			`{"hello":{"world":1234.99}}`,
			nil,
			map[string]interface{}{"hello.world": 1234.99},
		},
		{
			`{"hello":{"world":true}}`,
			nil,
			map[string]interface{}{"hello.world": true},
		},
		{
			`{"hello":{"world":null}}`,
			nil,
			map[string]interface{}{"hello.world": nil},
		},
		// empty slice
		{
			`{"hello":{"world":[]}}`,
			nil,
			map[string]interface{}{"hello.world": []interface{}{}},
		},
		// slice
		{
			`{"hello":{"world":["one","two"]}}`,
			nil,
			map[string]interface{}{
				"hello.world.0": "one",
				"hello.world.1": "two",
			},
		},
		// nested twice
		{
			`{"hello":{"world":{"again":"good morning"}}}`,
			nil,
			map[string]interface{}{"hello.world.again": "good morning"},
		},
		// multiple keys
		{
			`{
				"hello": {
					"lorem": {
						"ipsum":"again",
						"dolor":"sit"
					}
				},
				"world": {
					"lorem": {
						"ipsum":"again",
						"dolor":"sit"
					}
				}
			}`,
			nil,
			map[string]interface{}{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit"},
		},
		// empty object
		{
			`{"hello":{"empty":{"nested":{}}}}`,
			nil,
			map[string]interface{}{"hello.empty.nested": map[string]interface{}{}},
		},
		// custom delimiter
		{
			`{"hello":{"world":{"again":"good morning"}}}`,
			&Options{
				Delimiter: ":",
				MaxDepth:  20,
			},
			map[string]interface{}{"hello:world:again": "good morning"},
		},
		// custom depth
		{
			`{
				"hello": {
					"world": {
						"again": "good morning"
					}
				},
				"lorem": {
					"ipsum": {
						"dolor": "good evening"
					}
				}
			}
			`,
			&Options{
				MaxDepth:  2,
				Delimiter: ".",
			},
			map[string]interface{}{
				"hello.world": map[string]interface{}{"again": "good morning"},
				"lorem.ipsum": map[string]interface{}{"dolor": "good evening"},
			},
		},
		// custom safe = true
		{
			`{"hello":{"world":["one","two"]}}`,
			&Options{
				Safe:      true,
				Delimiter: ".",
			},
			map[string]interface{}{
				"hello.world": []interface{}{"one", "two"},
			},
		},
	}
	for i, test := range tests {
		var given interface{}
		err := json.Unmarshal([]byte(test.given), &given)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
		}
		got, err := Flatten(given.(map[string]interface{}), test.options)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestUnflatten(t *testing.T) {
	tests := []struct {
		flat    map[string]interface{}
		options *Options
		want    map[string]interface{}
	}{
		{
			map[string]interface{}{"hello": "world"},
			nil,
			map[string]interface{}{"hello": "world"},
		},
		// Key starts with upper case
		{
			map[string]interface{}{"Hello": "world"},
			nil,
			map[string]interface{}{"Hello": "world"},
		},
		{
			map[string]interface{}{"hello": 1234.56},
			nil,
			map[string]interface{}{"hello": 1234.56},
		},
		{
			map[string]interface{}{"hello": true},
			nil,
			map[string]interface{}{"hello": true},
		},
		// nested twice
		{
			map[string]interface{}{"hello.world.again": "good morning"},
			nil,
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		// multiple keys
		{
			map[string]interface{}{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit",
				"world":             map[string]interface{}{"greet": "hello"},
			},
			nil,
			map[string]interface{}{
				"hello": map[string]interface{}{
					"lorem": map[string]interface{}{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
				"world": map[string]interface{}{
					"greet": "hello",
					"lorem": map[string]interface{}{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
			},
		},
		// multiple keys - key starts with upper case
		{
			map[string]interface{}{
				"Hello.lorem.ipsum": "L1 upper",
				"hello.lorem.ipsum": "L1 lower",
				"hello.Lorem.dolor": "L2 upper",
				"hello.lorem.dolor": "L2 lower",
				"world.lorem.Ipsum": "L3 upper",
				"world.lorem.ipsum": "L3 lower",
				"world.lorem.dolor": "sit",
				"world": map[string]interface{}{
					"greet": "hello",
					"From":  "alice",
				},
			},
			nil,
			map[string]interface{}{
				"hello": map[string]interface{}{
					"lorem": map[string]interface{}{
						"ipsum": "L1 lower",
						"dolor": "L2 lower",
					},
					"Lorem": map[string]interface{}{"dolor": "L2 upper"},
				},
				"Hello": map[string]interface{}{
					"lorem": map[string]interface{}{"ipsum": "L1 upper"},
				},
				"world": map[string]interface{}{
					"greet": "hello",
					"From":  "alice",
					"lorem": map[string]interface{}{
						"ipsum": "L3 lower",
						"Ipsum": "L3 upper",
						"dolor": "sit",
					},
				},
			},
		},
		// nested objects do not clobber each other
		{
			map[string]interface{}{
				"foo.bar": map[string]interface{}{"t": 123},
				"foo":     map[string]interface{}{"k": 456},
			},
			nil,
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"t": 123,
					},
					"k": 456,
				},
			},
		},
		// custom delimiter
		{
			map[string]interface{}{
				"hello world again": "good morning",
			},
			&Options{
				Delimiter: " ",
			},
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		// do not overwrite
		{
			map[string]interface{}{
				"travis":           "true",
				"travis_build_dir": "/home/foo",
			},
			&Options{
				Delimiter: "_",
			},
			map[string]interface{}{
				"travis": "true",
			},
		},
		// keys with nil values
		{
			map[string]interface{}{
				"foo.bar": map[string]interface{}{"t": nil},
				"foo":     map[string]interface{}{"k": nil},
			},
			nil,
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"t": nil,
					},
					"k": nil,
				},
			},
		},
		// todo
		// overwrite true
		// {
		// 	map[string]interface{}{
		// 		"travis":           "true",
		// 		"travis_build_dir": "/home/foo",
		// 	},
		// 	Options{
		// 		Delimiter: "_",
		// 		Overwrite: true,
		// 	},
		// 	map[string]interface{}{
		// 		"travis": map[string]interface{}{
		// 			"build": map[string]interface{}{
		// 				"dir": "/home/foo",
		// 			},
		// 		},
		// 	},
		// },
	}
	for i, test := range tests {
		got, err := Unflatten(test.flat, test.options)
		if err != nil {
			t.Errorf("%d: failed to unflatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestFlattenPrefix(t *testing.T) {
	tests := []struct {
		given   string
		options *Options
		want    map[string]interface{}
	}{
		// test with different primitives
		// String: 'world',
		// Number: 1234.99,
		// Boolean: true,
		// null: null,
		{
			`{"hello": "world"}`,
			&Options{Prefix: "test", Delimiter: "."},
			map[string]interface{}{"test.hello": "world"},
		},
		{
			`{"hello": 1234.99}`,
			&Options{Prefix: "test", Delimiter: "_"},
			map[string]interface{}{"test_hello": 1234.99},
		},
		{
			`{"hello": true}`,
			&Options{Prefix: "test", Delimiter: "-"},
			map[string]interface{}{"test-hello": true},
		},
		{
			`{"hello":{"world":"good morning"}}`,
			&Options{Prefix: "test", Delimiter: "."},
			map[string]interface{}{"test.hello.world": "good morning"},
		},
		{
			`{"hello":{"world":1234.99}}`,
			&Options{Prefix: "test", Delimiter: "_"},
			map[string]interface{}{"test_hello_world": 1234.99},
		},
		{
			`{"hello":{"world":true}}`,
			&Options{Prefix: "test", Delimiter: "-"},
			map[string]interface{}{"test-hello-world": true},
		},
	}
	for i, test := range tests {
		var given interface{}
		err := json.Unmarshal([]byte(test.given), &given)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
		}
		got, err := Flatten(given.(map[string]interface{}), test.options)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestUnflattenPrefix(t *testing.T) {
	tests := []struct {
		flat    map[string]interface{}
		options *Options
		want    map[string]interface{}
	}{
		{
			map[string]interface{}{"test.hello": "world"},
			&Options{Prefix: "test", Delimiter: "."},
			map[string]interface{}{"hello": "world"},
		},
		{
			map[string]interface{}{"test_hello": 1234.56},
			&Options{Prefix: "test", Delimiter: "_"},
			map[string]interface{}{"hello": 1234.56},
		},
		{
			map[string]interface{}{"test-hello": true},
			&Options{Prefix: "test", Delimiter: "-"},
			map[string]interface{}{"hello": true},
		},
		// nested twice
		{
			map[string]interface{}{"test.hello.world.again": "good morning"},
			&Options{Prefix: "test", Delimiter: "."},
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		// custom delimiter
		{
			map[string]interface{}{
				"test hello world again": "good morning",
			},
			&Options{
				Prefix:    "test",
				Delimiter: " ",
			},
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
	}
	for i, test := range tests {
		got, err := Unflatten(test.flat, test.options)
		if err != nil {
			t.Errorf("%d: failed to unflatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

// --- Benchmarks ---

// Define package-level variables to store results and prevent compiler optimization.
var benchResult map[string]interface{}
var benchErr error

// prepareBenchmarkData creates a moderately complex nested map for benchmarking.
func prepareBenchmarkData() map[string]interface{} {
	return map[string]interface{}{
		"user": map[string]interface{}{ // L1
			"id":       12345,
			"name":     "John Doe",
			"email":    "john.doe@example.com",
			"isActive": true,
			"address": map[string]interface{}{ // L2
				"street":  "123 Main St",
				"city":    "Anytown",
				"zip":     "12345",
				"country": "USA",
			},
			"roles": []interface{}{ // L2
				"admin",
				"editor",
				map[string]interface{}{ // L3
					"permission": "read-only",
				},
			},
		},
		"settings": map[string]interface{}{ // L1
			"theme":    "dark",
			"fontSize": 14,
			"notify": map[string]interface{}{ // L2
				"email": true,
				"sms":   false,
			},
			"preferences": []interface{}{ // L2
				map[string]interface{}{ // L3
					"type":  "a",
					"value": true,
				},
				map[string]interface{}{ // L3
					"type":  "b",
					"value": false,
				},
			},
		},
		"metadata": map[string]interface{}{ // L1
			"createdAt": "2024-01-01T10:00:00Z",
			"updatedAt": "2024-04-12T15:30:00Z",
			"version":   3.1,
		},
	}
}

// prepareBenchmarkFlatData creates a flat version of the benchmark data.
func prepareBenchmarkFlatData() map[string]interface{} {
	// Flatten the nested data once to get the input for Unflatten benchmark
	data, _ := Flatten(prepareBenchmarkData(), nil)
	return data
}

func BenchmarkFlatten(b *testing.B) {
	nestedData := prepareBenchmarkData()
	opts := &Options{Delimiter: "."}
	var r map[string]interface{}
	var err error

	b.ReportAllocs()
	b.ResetTimer() // Reset timer after setup

	for i := 0; i < b.N; i++ {
		r, err = Flatten(nestedData, opts)
	}

	// Assign results to prevent optimization
	benchResult = r
	benchErr = err
}

func BenchmarkUnflatten(b *testing.B) {
	flatData := prepareBenchmarkFlatData() // Prepare flat data outside the loop
	opts := &Options{Delimiter: "."}
	var r map[string]interface{}
	var err error

	b.ReportAllocs()
	b.ResetTimer() // Reset timer after setup

	for i := 0; i < b.N; i++ {
		r, err = Unflatten(flatData, opts)
	}

	// Assign results to prevent optimization
	benchResult = r
	benchErr = err
}
