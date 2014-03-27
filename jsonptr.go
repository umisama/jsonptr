package jsonptr

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

func Find(src []byte, path string) (obj []byte, err error) {
	var srcjson interface{}
	err = json.Unmarshal(src, &srcjson)
	if err != nil {
		return
	}

	splitedpath := pathProcessor(path)
	obj_raw, err := find(srcjson, splitedpath, 0)

	obj, err = json.Marshal(obj_raw)
	return
}

func pathProcessor(path string) (processed []string) {
	if len(path) == 0 {
		return
	}

	if path[0] == '#' {
		processed = pathProcessorSub(path[1:], true)
	} else {
		processed = pathProcessorSub(path, false)
	}

	return
}

func pathProcessorSub(path string, uriencode bool) (processed []string) {
	processed = make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if v == "" {
			continue
		}

		reper := strings.NewReplacer("~0", "~", "~1", "/", "\\\\", "\\", "\\\"", "\"")
		v = reper.Replace(v)

		if uriencode {
			v, _ = url.QueryUnescape(v)
		}

		processed = append(processed, v)
	}

	// it pointing empty string attr, if path end with "/".
	if strings.HasSuffix(path, "/") {
		processed = append(processed, "")
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
