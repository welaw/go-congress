package utils

import (
	"io/ioutil"
	"net/http"
)

func DownloadData(p string) (dat []byte, err error) {
	res, err := http.Get(p)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dat, err = ioutil.ReadAll(res.Body)
	return dat, err
}
