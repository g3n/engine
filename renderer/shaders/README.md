This directory contains GLSL shaders used by the engine
-------------------------------------------------------

If any shader in this directory or include 'chunk' in the
"include" subdirectory is modified or a new shader or chunk
is added or removed it is necessary to execute:

>go generate

in this directory to update the "sources.go" file.
It will invoke the "g3nshaders" command which will read
the shaders and include files and generate the "sources.go" file.

To install "g3nshaders" change to the "tools/g3nshaders" directory
from the engine "root" and execute: "go install".

