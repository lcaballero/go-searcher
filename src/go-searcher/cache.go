package searcher

import (
	"io/ioutil"
	"os"
	"errors"
	"path/filepath"
	"encoding/json"
)

type CacheFile struct {
	Name string
	Content string
	Basename string
	Extension string
	Dir string
	IsDirectory bool
	Error error
}

type Index struct {
	Root string
	Files []*CacheFile
}

type SearchOptions struct {
	Cwd ToCwd
}

type Searcher struct {
	cwd ToCwd
}

type ToCwd func(path string) string

func defaultCwd(f string) string {
	return filepath.Join(".", f)
}

func (index *Index) ToJson() []byte {
	bytes,_ := json.MarshalIndent(index, "", "  ")
	return bytes
}

func (index *Index) ToFile(filename string) error {
	//	ioutil.WriteFile(filename, index.ToJson(), os.)
	//	os.Create(filename)
	return ioutil.WriteFile(filename, index.ToJson(), 0644)
	//	r   w   x
	//	110 010 010
	//	6   4   4
}

func New(opts SearchOptions) *Searcher {
	var cwd ToCwd = defaultCwd
	if opts.Cwd != nil {
		cwd = opts.Cwd
	}
	return &Searcher{
		cwd:cwd,
	}
}

func (index *Index) Names() []string {
	names := make([]string, len(index.Files))
	for i,n := range index.Files {
		names[i] = n.Name
	}
	return names
}

func newCacheFile(root, name string, content string, isDir bool) *CacheFile {
	return &CacheFile{
		Name:name,
		Content:content,
		Basename:filepath.Base(name),
		Extension:filepath.Ext(name),
		Dir:root,
		IsDirectory:isDir,
	}
}

func readDir(root string) ([]os.FileInfo, error){
	var err error = nil
	var rootInfo os.FileInfo

	if rootInfo,err = os.Stat(root); os.IsNotExist(err) {
		return nil, err
	}

	if !rootInfo.IsDir() {
		return nil, errors.New("Expected root to be directory")
	}

	dirs,err := ioutil.ReadDir(root)

	if err != nil {
		return nil, err
	}

	return dirs, err
}

func extendIndex(root string, index *Index) {
	dirs, err := readDir(root)
	if err != nil {
		return
	}
	for _, info := range dirs {
		name := info.Name()
		var cacheFile *CacheFile = nil
		fullPath := filepath.Join(root, name)

		if info.IsDir() {
			cacheFile = newCacheFile(root, name, "", false)
			extendIndex(fullPath, index)
		} else {
			bytes,_ := ioutil.ReadFile(fullPath)
			content := string(bytes)
			cacheFile = newCacheFile(root, name, content, true)
		}
		index.Files = append(index.Files, cacheFile)
	}
}

func (s *Searcher) BuildIndex(rel string) (*Index, error) {
	root := s.cwd(rel)
	dirs,err := readDir(root)

	if err != nil {
		return nil, err
	}

	index := &Index{
		Root:root,
		Files: make([]*CacheFile, 0),
	}

	for _, info := range dirs {
		name := info.Name()
		var cacheFile *CacheFile = nil
		fullPath := filepath.Join(root, name)

		if info.IsDir() {
			cacheFile = newCacheFile(root, name, "", true)
			extendIndex(fullPath, index)
		} else {
			bytes,_ := ioutil.ReadFile(fullPath)
			content := string(bytes)
			cacheFile = newCacheFile(root, name, content, false)
		}
		index.Files = append(index.Files, cacheFile)
	}

	return index, nil
}
