// Copyright 2014 Music Glue. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
/*
Package oak is an implementation of a Trie-style data structure
implemented with mutexes around the nodes to ensure coherent and
safe data access and updating from multiple concurrent accessors.

Sample use cases would include a lookup table for routes in a web
process.
*/
package oak
