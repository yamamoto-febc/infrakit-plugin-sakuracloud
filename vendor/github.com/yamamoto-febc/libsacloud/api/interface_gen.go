package api

/************************************************
  generated by IDE. for [InterfaceAPI]
************************************************/

import (
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

/************************************************
   To support influent interface for Find()
************************************************/

func (api *InterfaceAPI) Reset() *InterfaceAPI {
	api.reset()
	return api
}

func (api *InterfaceAPI) Offset(offset int) *InterfaceAPI {
	api.offset(offset)
	return api
}

func (api *InterfaceAPI) Limit(limit int) *InterfaceAPI {
	api.limit(limit)
	return api
}

func (api *InterfaceAPI) Include(key string) *InterfaceAPI {
	api.include(key)
	return api
}

func (api *InterfaceAPI) Exclude(key string) *InterfaceAPI {
	api.exclude(key)
	return api
}

func (api *InterfaceAPI) FilterBy(key string, value interface{}) *InterfaceAPI {
	api.filterBy(key, value, false)
	return api
}

// func (api *InterfaceAPI) FilterMultiBy(key string, value interface{}) *InterfaceAPI {
// 	api.filterBy(key, value, true)
// 	return api
// }

func (api *InterfaceAPI) WithNameLike(name string) *InterfaceAPI {
	return api.FilterBy("Name", name)
}

func (api *InterfaceAPI) WithTag(tag string) *InterfaceAPI {
	return api.FilterBy("Tags.Name", tag)
}
func (api *InterfaceAPI) WithTags(tags []string) *InterfaceAPI {
	return api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *InterfaceAPI) WithSizeGib(size int) *InterfaceAPI {
// 	api.FilterBy("SizeMB", size*1024)
// 	return api
// }

// func (api *InterfaceAPI) WithSharedScope() *InterfaceAPI {
// 	api.FilterBy("Scope", "shared")
// 	return api
// }

// func (api *InterfaceAPI) WithUserScope() *InterfaceAPI {
// 	api.FilterBy("Scope", "user")
// 	return api
// }

func (api *InterfaceAPI) SortBy(key string, reverse bool) *InterfaceAPI {
	api.sortBy(key, reverse)
	return api
}

func (api *InterfaceAPI) SortByName(reverse bool) *InterfaceAPI {
	api.sortByName(reverse)
	return api
}

// func (api *InterfaceAPI) SortBySize(reverse bool) *InterfaceAPI {
// 	api.sortBy("SizeMB", reverse)
// 	return api
// }

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

func (api *InterfaceAPI) New() *sacloud.Interface {
	return &sacloud.Interface{}
}

func (api *InterfaceAPI) Create(value *sacloud.Interface) (*sacloud.Interface, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.create(api.createRequest(value), res)
	})
}

func (api *InterfaceAPI) Read(id int64) (*sacloud.Interface, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.read(id, nil, res)
	})
}

func (api *InterfaceAPI) Update(id int64, value *sacloud.Interface) (*sacloud.Interface, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *InterfaceAPI) Delete(id int64) (*sacloud.Interface, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.delete(id, nil, res)
	})
}

/************************************************
  Inner functions
************************************************/

func (api *InterfaceAPI) setStateValue(setFunc func(*sacloud.Request)) *InterfaceAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

func (api *InterfaceAPI) request(f func(*sacloud.Response) error) (*sacloud.Interface, error) {
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.Interface, nil
}

func (api *InterfaceAPI) createRequest(value *sacloud.Interface) *sacloud.Request {
	req := &sacloud.Request{}
	req.Interface = value
	return req
}
