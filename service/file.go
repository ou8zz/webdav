package service

import (
  "fmt"
  "io/ioutil"
)

const layoutDateTime = "2006-01-02 15:04:05"

type (
  Item struct {
    Name    string `json:"name"`
    ModTime string `json:"modTime"`
    Size    int64  `json:"size"`
    IsDir   bool   `json:"isDir"`
  }
)

func GetFiles(root string) []*Item {
  result := make([]*Item, 0)
  f, err := ioutil.ReadDir(root)
  if err != nil {
    fmt.Printf("GetFiles files %+v, err:%v \n", f, err)
  }
  for _, v := range f {
    vo := &Item{}
    vo.Name = v.Name()
    vo.Size = v.Size()
    vo.IsDir = v.IsDir()
    vo.ModTime = v.ModTime().Format(layoutDateTime)
    result = append(result, vo)
  }
  return result
}
