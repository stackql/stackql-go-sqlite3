# stackql-go-sqlite3

> forked from [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3)

Embedded Golang SQLite Distribution for StackQL

## Adding StackQL Extension Functions for SQLite3

The `feature/stackql-ext-fns` branch of `stackql/stackql-go-sqlite3` (this repo) contains additional files and modifications to `sqlite3-binding.c` to add required extension functions to StackQL (using the default embedded `sqlite` backend).  Some of these functions may not be available using a `postgres` backend for StackQL.

### Source for StackQL Extension Functions

SQLite extension functions for StackQL are written in `C` and mastered in [`stackql/sqlite-ext-functions`](https://github.com/stackql/sqlite-ext-functions), user documentation for these functions can be found in this repo.

### Including StackQL Extension Functions in `stackql-go-sqlite3`

StackQL extension functions are included with the preprocessor directive `SQLITE_ENABLE_STACKQL`, which is defined in `sqlite3_opt_stackql.go`, which needs to be added to the root of this repo, the file contains this...

```golang
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//go:build sqlite_stackql || stackql
// +build sqlite_stackql stackql

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_ENABLE_STACKQL
#cgo LDFLAGS: -lm
*/
import "C"
```

Additional header files must be added to the root to support custom functions, these header functions and their associated program files are stored with the various functions in [`stackql/sqlite-ext-functions`](https://github.com/stackql/sqlite-ext-functions).  

Additions to the `sqlite3-binding.c` megalith amalgamation file, in the appropriate places, as shown here:

```c
/*
** Forward declarations of external module initializer functions
** for modules that need them.
*/
...other existing stuff
// CUSTOM_EXTENSIONS
// 1. Add custom extension initializer function declarations here...
#ifdef SQLITE_ENABLE_STACKQL
SQLITE_PRIVATE int sqlite3_jsonequal_init(sqlite3*, char**, const sqlite3_api_routines*);
SQLITE_PRIVATE int sqlite3_regexp_init(sqlite3*, char**, const sqlite3_api_routines*);
SQLITE_PRIVATE int sqlite3_splitpart_init(sqlite3*, char**, const sqlite3_api_routines*);
#endif
// End CUSTOM_EXTENSIONS
/*
** An array of pointers to extension initializer functions for
** built-in extensions.
*/
static int (*const sqlite3BuiltinExtensions[])(sqlite3*) = {
// CUSTOM_EXTENSIONS
// 2. Include custom extension initializer functions here...
#ifdef SQLITE_ENABLE_STACKQL
  (int(*)(sqlite3*))sqlite3_jsonequal_init,
  (int(*)(sqlite3*))sqlite3_regexp_init,
  (int(*)(sqlite3*))sqlite3_splitpart_init,
#endif
// End CUSTOM_EXTENSIONS  
...other existing stuff
};
...other existing stuff
#endif /* SQLITE_USER_AUTHENTICATION */

// CUSTOM_EXTENSIONS
// 3. Add custom extension initializer functions here...
#ifdef SQLITE_ENABLE_STACKQL

// ...custom function code here for all functions...

#endif // #ifdef SQLITE_ENABLE_STACKQL
// End CUSTOM_EXTENSIONS
```

### Building `stackql-go-sqlite3` 

To test compile `stackql-go-sqlite3`, run the following command:

```bash
go build --tags "sqlite_stackql" -o /dev/null
```

## Deployment to StackQL

`stackql-go-sqlite3` is pulled in via the CI processes for [`stackql/stackql`](https://github.com/stackql/stackql) or [`stackql/stackql-devel`](https://github.com/stackql/stackql-devel) via git tags.  Tags should follow the convention below:  

```bash
v{sqlite_major_version}.{sqlite_minor_version}.{sqlite_patch_version}-stackql{YYYYMMDD}[-{seq_no}]
```
`seqno` is only used if more than one tag is pushed on a given day.  An example of the steps to add and push tags are shown here:  

```bash
git add .
git commit -m "function updates"
git push origin feature/stackql-ext-fns
git tag v3.45.1-stackql20240708
git push origin v3.45.1-stackql20240708
```

### updating `go.mod` in `stackql`

In [`stackql/stackql`](https://github.com/stackql/stackql) or [`stackql/stackql-devel`](https://github.com/stackql/stackql-devel), update `go.mod` to reflect the current tag for `stackql-go-sqlite3`, for example...

```go
// ...
require (
// ...other stuff
	github.com/stackql/stackql-go-sqlite3 v3.45.1-stackql20240708
// ...other stuff
)
// ...
```

> Ensure tests are added and pass for any additional functions included.