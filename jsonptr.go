package jsonptr

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

var (
	replacer              = strings.NewReplacer("~0", "~", "~1", "/", "\\\\", "\\", "\\\"", "\"")
	ErrInvalidJsonPointer = errors.New("jsonptr: invalid json pointer")
	ErrInvalidJsonString  = errors.New("jsonptr: invalid json object")
)

func Find(src []byte, path string) (obj []byte, err error) {
	var srcjson interface{}
	err = json.Unmarshal(src, &srcjson)
	if err != nil {
		err = ErrInvalidJsonString
		return
	}

	var obj_raw interface{}
	if path == "" || path == "#" {
		obj_raw = srcjson
	} else {
		var spath []string
		spath, err = pathProcessor(path)
		obj_raw, err = find(srcjson, spath, 0)
	}

	obj, err = json.Marshal(obj_raw)
	if err != nil {
		err = ErrInvalidJsonPointer
		return
	}

	return
}

func pathProcessor(path string) (processed []string, err error) {
	if len(path) == 0 {
		err = ErrInvalidJsonPointer
		return
	}

	switch {
	case strings.HasPrefix(path, "#/"):
		processed, err = pathProcessorSub(path[2:], true)
	case strings.HasPrefix(path, "/"):
		processed, err = pathProcessorSub(path[1:], false)
	default:
		err = ErrInvalidJsonPointer
	}

	return
}

func pathProcessorSub(path string, uriencode bool) (processed []string, err error) {
	processed = make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		v = replacer.Replace(v)
		if uriencode {
			v, err = url.QueryUnescape(v)
		}

		processed = append(processed, v)
	}

	if len(processed) == 0 {
		processed = []string{""}
	}
	return
}

func find(src interface{}, path []string, level int) (obj interface{}, err error) {
	if len(path) == level {
		obj = src
		return
	}

	i, name := 0, path[level]
	switch t := src.(type) {
	case map[string]interface{}:
		obj, err = find(t[name], path, level+1)
	case []interface{}:
		i, err = strconv.Atoi(name)
		if err != nil {
			return
		}
		obj, err = find(t[i], path, level+1)
	}

	return
}
