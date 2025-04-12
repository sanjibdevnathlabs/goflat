# goflat

[![Go Build](https://github.com/sanjibdevnathlabs/goflat/actions/workflows/test.yml/badge.svg)](https://github.com/sanjibdevnathlabs/goflat/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/sanjibdevnathlabs/goflat/graph/badge.svg)](https://codecov.io/gh/sanjibdevnathlabs/goflat)
[![Go Report Card](https://goreportcard.com/badge/github.com/sanjibdevnathlabs/goflat)](https://goreportcard.com/report/github.com/sanjibdevnathlabs/goflat)
[![Go Reference](https://pkg.go.dev/badge/github.com/sanjibdevnathlabs/goflat.svg)](https://pkg.go.dev/github.com/sanjibdevnathlabs/goflat)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This project has been forked from [https://github.com/nqd/flat](https://github.com/nqd/flat)

Take a golang map and flatten it or unflatten a map with delimited key.

## Method

### Flatten

Flatten given map, returns a map one level deep.

```{go}
in := map[string]interface{}{
    "a": "b",
    "c": map[string]interface{}{
        "d": "e",
        "f": "g",
    },
    "z": [2, 1.4567],
}

out, err := flat.Flatten(in, nil)
// out = map[string]interface{}{
//     "a": "b",
//     "c.d": "e",
//     "c.f": "g",
//     "z.0": 2,
//     "z.1": 1.4567,
// }
```

### Unflatten

Since there is flatten, flat should have unfatten.

```{go}
in := map[string]interface{}{
    "foo.bar": map[string]interface{}{"t": 123},
    "foo":     map[string]interface{}{"k": 456},
}

out, err := flat.Unflatten(in, nil)
// out = map[string]interface{}{
//     "foo": map[string]interface{}{
//         "bar": map[string]interface{}{
//             "t": 123,
//         },
//         "k": 456,
//     },
// }
```

## Options

### Delimiter

Use a custom delimiter for flattening/unflattening your objects. Default value is `.`.

```{go}
in := map[string]interface{}{
   "hello": map[string]interface{}{
       "world": map[string]interface{}{
           "again": "good morning",
        }
    },
}

out, err := flat.Flatten(in, &flat.Options{
    Delimiter: ":",
})
// out = map[string]interface{}{
//     "hello:world:again": "good morning",
// }
```

### Safe

<!-- When Safe is true, both fatten and unflatten will preserve arrays and their contents. Default Safe value is `false`. -->
When Safe is true, fatten will preserve arrays and their contents. Default Safe value is `false`.

```{go}
in := map[string]interface{}{
    "hello": map[string]interface{}{
        "world": []interface{}{
            "one",
            "two",
        }
   },
}

out, err := flat.Flatten(in, &flat.Options{
    Delimiter: ".",
    Safe:      true,
})
// out = map[string]interface{}{
//     "hello.world": []interface{}{"one", "two"},
// }
```

<!-- Example of Unflatten goes here -->

### MaxDepth

MaxDepth is the maximum number of nested objects to flatten. MaxDepth can be any integer number. MaxDepth = 0 means no limit.

Default MaxDepth value is `0`.

```{go}
in := map[string]interface{}{
    "hello": map[string]interface{}{
        "world": []interface{}{
            "again": "good morning",
        }
   },
}

out, err := flat.Flatten(in, &flat.Options{
    Delimiter: ".",
    MaxDepth:  2,
})
// out = map[string]interface{}{
//     "hello.world": map[string]interface{}{"again": "good morning"},
// }
```

## Todos

- [ ] Safe option for Unflatten
- [ ] Overwrite

## Known Issues

- If slice is flattened, there are issues with unflattening it.
