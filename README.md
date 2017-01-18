In-Memory file system for golang
================================

filesystem-internal implements an adapter for in-memory files for the golang
filesystem integration (see the filesystem package).

The file system is registered as "internal://" and only respects the path,
not the host (i.e. internal://one/path1 is the same as internal://two/path1
but internal:///one/path1 is different from internal:///two/path1). The file
system registers itself automatically upon loading the module.

Usage
=====

Just load the internal file system module:

    import (
        "github.com/childoftehuniverse/filesystem"
        _ "github.com/childoftehuniverse/filesystem-internal"
    )

Then you can just use the regular filesystem API as documented. See
<https://github.com/childoftheuniverse/filesystem/> for details.

Alternatively, you can just use the anonymous in-memory file implementation:

    import (
        internal "github.com/childoftehuniverse/filesystem-internal"
    )

[…]

    var file = internal.NewAnonymousFile()
    file.Write(…)
