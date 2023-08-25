# Oh My Glob

A library for file path globbing in Go, with specific support for the `**` construct (which is not handled by `filepath.Glob`.)

```go
path := "foo/bar/baz.txt"
my_glob := = oh_my_glob.Compile("foo/**/*.txt")
if my_glob.Match() {
        // et cetera
```

## Usage

Globs in `oh_my_glob` are represented using a `Glob` struct, which contained a compiled representation for fast matching.

```go
// A struct representing a compiled glob
type Glob
```

You can compile a path—represented as a string—to a `Glob` struct using the `Compile` function.

```go
func Compile(glob string) Glob
```

Once you have a `Glob`, you can match it against an explicit path using the `Match` method.

```go
func (g *Glob) Match(path string) bool
```

Additionally, every `Glob` remembers the string it was compiled from, which might be useful for categorization or debugging:

```go
func (g *Glob) Original() string
```

## Understood features

Most string literals are understood by `oh_my_glob` as matching only themselves.

- The asterisk `*` matches any string of any length not including a slash. For example, the glob `x*` will match the paths `"x"`, `"xa"`, and `"xaaa"`, but it will not match `"x/a"`.
- The double asterisk `**` matches any number of relative directories, including zero. For example, the glob `a/**/b` will match the paths `"a/x/y/b"`, `"a/x/b"`, and even `"a/b"`. (Note that this means a glob can have more forward slashes than a path it recognizes!) The glob `a/**` will match anything contained within `a/` as well as the path `"a"` by itself, and the glob `**/b` will match anything named `b` within any directories, including the path `"b"` itself.
