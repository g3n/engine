package assets

// To generate file with fonts binary data install "go-bindata" from:
// https://github.com/jteeuwen/go-bindata
// >go get -u github.com/jteeuwen/go-bindata/...

//go:generate go-bindata -o data.go -pkg assets fonts
//go:generate g3nicodes -pkg icon icon/codepoints icon/icodes.go
