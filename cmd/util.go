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

// unmarshalOverrides returns overrides converted from a JSON string or the contents of a JSON file.
func unmarshalOverrides(override string) ([]*api.OverrideSpec, error) {
	if strings.HasPrefix(override, "file://") {
		bytes, err := ioutil.ReadFile(override[7:])
		if err != nil {
			return nil, err
		}
		override = string(bytes)
	}
	if strings.HasPrefix(override, "{") {
		var o api.OverrideSpec
		if err := protojson.Unmarshal([]byte(override), &o); err != nil {
			return nil, err
		}
		return []*api.OverrideSpec{&o}, nil
	}
	if strings.HasPrefix(override, "[") {
		var os api.OverrideSpecs
		if err := protojson.Unmarshal([]byte(override), &os); err != nil {
			return nil, err
		}
		return os.Specs, nil
	}
	return nil, nil
}
