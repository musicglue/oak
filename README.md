oak
===

[![Build Status](https://travis-ci.org/musicglue/oak.svg?branch=master)](https://travis-ci.org/musicglue/oak)

Oak is a trie-based data structure designed for mapping nested strings to disparate interface values.

There'll be a readme here soon, but the best place to go to work out how this works is [GoDoc](http://godoc.org/github.com/musicglue/oak).

## Example

For now though, this will map a slice of strings (think URL segments) to a value (interface{}) that is stored against them.

```go
package main

import(
  "github.com/musicglue/oak"
  "fmt"
)

func main() {
  tree := oak.NewBranch()
  tree.Set([]string{}, "Home")
  tree.Set([]string{"pages", "about-us"}, "About Us")

  home, ok := tree.Get([]string{})
  if ok {
    fmt.Println("Home resolves to:", home)
  }

  about_us, ok := tree.Get([]string{"pages", "about-us"})
  if ok {
    fmt.Println("About Us resolves to:", about_us)
  }

  fallover, ok := tree.Match([]string{"pages"})
  if ok {
    fmt.Println("Calling a non-extant key with Match returns the deepest extant value", fallover)
  }
}
```
