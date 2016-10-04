package api

/************************************************
  generated by IDE. for [CDROMAPI]
************************************************/

import (
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

/************************************************
   To support influent interface for Find()
************************************************/

func (api *CDROMAPI) Reset() *CDROMAPI {
	api.reset()
	return api
}

func (api *CDROMAPI) Offset(offset int) *CDROMAPI {
	api.offset(offset)
	return api
}

func (api *CDROMAPI) Limit(limit int) *CDROMAPI {
	api.limit(limit)
	return api
}

func (api *CDROMAPI) Include(key string) *CDROMAPI {
	api.include(key)
	return api
}

func (api *CDROMAPI) Exclude(key string) *CDROMAPI {
	api.exclude(key)
	return api
}

func (api *CDROMAPI) FilterBy(key string, value interface{}) *CDROMAPI {
	api.filterBy(key, value, false)
	return api
}

// func (api *CDROMAPI) FilterMultiBy(key string, value interface{}) *CDROMAPI {
// 	api.filterBy(key, value, true)
// 	return api
// }

func (api *CDROMAPI) WithNameLike(name string) *CDROMAPI {
	return api.FilterBy("Name", name)
}

func (api *CDROMAPI) WithTag(tag string) *CDROMAPI {
	return api.FilterBy("Tags.Name", tag)
}
func (api *CDROMAPI) WithTags(tags []string) *CDROMAPI {
	return api.FilterBy("Tags.Name", []interface{}{tags})
}

func (api *CDROMAPI) WithSizeGib(size int) *CDROMAPI {
	api.FilterBy("SizeMB", size*1024)
	return api
}

func (api *CDROMAPI) WithSharedScope() *CDROMAPI {
	api.FilterBy("Scope", "shared")
	return api
}

func (api *CDROMAPI) WithUserScope() *CDROMAPI {
	api.FilterBy("Scope", "user")
	return api
}

func (api *CDROMAPI) SortBy(key string, reverse bool) *CDROMAPI {
	api.sortBy(key, reverse)
	return api
}

func (api *CDROMAPI) SortByName(reverse bool) *CDROMAPI {
	api.sortByName(reverse)
	return api
}

func (api *CDROMAPI) SortBySize(reverse bool) *CDROMAPI {
	api.sortBy("SizeMB", reverse)
	return api
}

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

func (api *CDROMAPI) New() *sacloud.CDROM {
	return &sacloud.CDROM{
		TagsType: &sacloud.TagsType{},
	}
}

//func (api *CDROMAPI) Create(value *sacloud.CDROM) (*sacloud.CDROM, error) {
//	return api.request(func(res *sacloud.Response) error {
//		return api.create(api.createRequest(value), res)
//	})
//}

func (api *CDROMAPI) Read(id int64) (*sacloud.CDROM, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.read(id, nil, res)
	})
}

func (api *CDROMAPI) Update(id int64, value *sacloud.CDROM) (*sacloud.CDROM, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *CDROMAPI) Delete(id int64) (*sacloud.CDROM, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.delete(id, nil, res)
	})
}

/************************************************
  Inner functions
************************************************/

func (api *CDROMAPI) setStateValue(setFunc func(*sacloud.Request)) *CDROMAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

func (api *CDROMAPI) request(f func(*sacloud.Response) error) (*sacloud.CDROM, error) {
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.CDROM, nil
}

func (api *CDROMAPI) createRequest(value *sacloud.CDROM) *sacloud.Request {
	req := &sacloud.Request{}
	req.CDROM = value
	return req
}
