package memory

func NewFileSystem(dir string) *FileSystem {
	return &FileSystem{dir: dir}
}

type FileSystem struct {
	dir string
}
