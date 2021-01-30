package utils

import "errors"

func GetMapPath(in interface{}, paths ...string) (interface{}, error) {
	if len(paths) == 0 {
		return in, nil
	}

	inMap, ok := in.(map[string]interface{})
	if !ok {
		return nil, errors.New("Can't get sub path " + paths[0] + ", because node is not a map")
	}

	tmp, ok := inMap[paths[0]]
	if !ok {
		return nil, errors.New("Do not has sub path" + paths[0])
	}

	p := paths[1:]
	return GetMapPath(tmp, p...)
}
