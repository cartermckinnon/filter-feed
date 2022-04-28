package cmd

import (
	"io/ioutil"
	"strings"

	api "github.com/cartermckinnon/filter-feed/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// unmarshalFilters returns filters converted from a JSON string or the contents of a JSON file.
func unmarshalFilters(filter string) ([]*api.FilterSpec, error) {
	if strings.HasPrefix(filter, "file://") {
		bytes, err := ioutil.ReadFile(filter[7:])
		if err != nil {
			return nil, err
		}
		filter = string(bytes)
	}
	if strings.HasPrefix(filter, "{") {
		var f api.FilterSpec
		if err := protojson.Unmarshal([]byte(filter), &f); err != nil {
			return nil, err
		}
		return []*api.FilterSpec{&f}, nil
	}
	if strings.HasPrefix(filter, "[") {
		var fs api.FilterSpecs
		if err := protojson.Unmarshal([]byte(filter), &fs); err != nil {
			return nil, err
		}
		return fs.Specs, nil
	}
	return nil, nil
}
