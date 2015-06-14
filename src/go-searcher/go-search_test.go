package searcher

import (
	"testing"
	"github.com/Workiva/go-datastructures/set"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"os"
)


func NewCwdSearcher() *Searcher {
	return New(SearchOptions{
		Cwd:func (f string) string {
			cwd := os.Getenv("TESTING_PWD")
			return filepath.Join(cwd, f)
		},
	})
}

func TestEmptyDirectory(t *testing.T) {

	Convey("Should entries in empty sub-directories", t, func() {
		index,err := NewCwdSearcher().BuildIndex("files/sources/s3")
		So(err, ShouldBeNil)
		So(index.Files, ShouldNotBeEmpty)
		So(len(index.Files), ShouldEqual, 4)

		set := set.New()
		for _,name := range index.Names() {
			set.Add(name)
		}

		So(set.All("b", "b.txt", "a", "a.txt"), ShouldBeTrue)
	})

	Convey("Should not find any entries in empty directory", t, func() {
		index,_ := NewCwdSearcher().BuildIndex("files/sources/empty")
		So(index.Files, ShouldBeEmpty)
	})

	Convey("Should find that not-dir/ doesn't exist", t, func() {
		_, err := NewCwdSearcher().BuildIndex("files/sources/not-dir")
		So(err, ShouldNotBeNil)
	})


	Convey("Should find empty entry from s2/empty-dir/", t, func() {
		index, err := NewCwdSearcher().BuildIndex("files/sources/s2/")
		So(err, ShouldBeNil)
		info := index.Files[0]
		So(info.Name, ShouldEqual, "empty-dir")
		So(info.Basename, ShouldEqual, "empty-dir")
		So(info.Extension, ShouldEqual, "")
		So(info.Dir, ShouldEndWith, "files/sources/s2")
		So(info.HasContent, ShouldBeFalse)
		So(info.IsDirectory, ShouldBeTrue)
	})

	Convey("Should find 1 entry in empty s1/", t, func() {
		index,_ := NewCwdSearcher().BuildIndex("files/sources/s1")
		So(index.Files, ShouldNotBeEmpty)
		So(len(index.Files), ShouldEqual, 1)
	})

	Convey("Should find entry with correct fields in s1/", t, func() {
		index, _ := NewCwdSearcher().BuildIndex("files/sources/s1")
		info := index.Files[0]
		So(info.Name, ShouldEqual, "empty-1.txt")
		So(info.Basename, ShouldEqual, "empty-1.txt")
		So(info.Extension, ShouldEqual, ".txt")
		So(info.Dir, ShouldEndWith, "files/sources/s1")
		So(info.HasContent, ShouldBeTrue)
		So(info.IsDirectory, ShouldBeFalse)
		So(info.Error, ShouldBeNil)
	})
}