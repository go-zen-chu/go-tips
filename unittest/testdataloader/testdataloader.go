package testdataloader

import (
	"os"
	"path/filepath"
)

// data 内に相対パスでテストデータが入るが、テストの仕方次第で data ではなく、わかりやすいデータ構造に入れておくのもあり
type testDataLoader struct {
	data map[string]string
}

func NewTestDataLoader() *testDataLoader {
	tdl := &testDataLoader{
		data: make(map[string]string),
	}
	if err := tdl.loadTestData(); err != nil {
		panic(err)
	}
	return tdl
}

func (t *testDataLoader) loadTestData() error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	testDir := filepath.Join(curDir, "testdata")
	return filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(testDir, path)
		if err != nil {
			return err
		}
		bt, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		t.data[relPath] = string(bt)
		return nil
	})
}
