// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//go:build sqlite_stackql || stackql
// +build sqlite_stackql stackql

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_ENABLE_STACKQL -DSQLITE_ENABLE_MATH_FUNCTIONS
#cgo LDFLAGS: -lm
*/
import "C"
