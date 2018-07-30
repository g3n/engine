package assets

// To generate file with fonts binary data install "go-bindata" from:
// https://github.com/go-bindata/go-bindata
// > go get -u github.com/go-bindata/go-bindata/...

//go:generate go-bindata -o data.go -pkg assets fonts cursors
//go:generate g3nicodes -pkg icon icon/codepoints icon/icodes.go
