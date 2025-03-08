/* For license and copyright information please see the LEGAL file in the code repository */

package service

import (
	"libgo/log"
	"libgo/protocol"
)

// SS is the same as the Services.
// Use this type when embed in other struct to solve field & method same name problem(Services struct and Services() method) to satisfy interfaces.
type SS = Services

// Services store all application service
type Services struct {
	poolByRegisterTime []protocol.Service
	poolByID           map[protocol.MediaTypeID]protocol.Service
	poolByMediaType    map[string]protocol.Service
	poolByURIPath      map[string]protocol.Service
}

// Init use to initialize
func (ss *Services) Init() (err protocol.Error) {
	const poolSizes = 512
	// TODO::: decide about poolSize by hardware

	ss.poolByRegisterTime = make([]protocol.Service, poolSizes)
	ss.poolByID = make(map[protocol.MediaTypeID]protocol.Service, poolSizes)
	ss.poolByURIPath = make(map[string]protocol.Service, poolSizes)
	ss.poolByMediaType = make(map[string]protocol.Service, poolSizes)
	return
}
func (ss *Services) Reinit() (err protocol.Error) {
	// TODO:::
	// for _, s := range ss.poolByURIPath {
	// 	s.Reinit()
	// }
	return
}
func (ss *Services) Deinit() (err protocol.Error) {
	for _, s := range ss.poolByURIPath {
		err = s.Deinit()
		// TODO::: easily return if occur any error??
		if err != nil {
			return
		}
	}
	return
}

// RegisterService use to register application services.
// Due to minimize performance impact, This method isn't safe to use concurrently and
// must register all service before use GetService methods.
//
//libgo:impl libgo/protocol.Services
func (ss *Services) RegisterService(s protocol.Service) (err protocol.Error) {
	if s.ID() == 0 && s.URI() == "" {
		err = &ErrServiceNotProvideIdentifier
		return
	}

	ss.registerServiceByMediaType(s)
	ss.registerServiceByURI(s)
	ss.poolByRegisterTime = append(ss.poolByRegisterTime, s)
	return
}

func (ss *Services) registerServiceByMediaType(s protocol.Service) (err protocol.Error) {
	var serviceID = s.ID()
	var exitingServiceByID, _ = ss.GetServiceByID(serviceID)
	if exitingServiceByID != nil {
		err = &ErrServiceDuplicateIdentifier
		log.Fatal(s, "ID associated for '"+s.MediaType()+"' Used before for other service and not legal to reuse same ID for other services\n"+
			"	Exiting service MediaType is: "+exitingServiceByID.MediaType())
	} else {
		ss.poolByID[serviceID] = s
	}

	var serviceMediaType = s.MediaType()
	var exitingServiceByMediaType, _ = ss.GetServiceByMediaType(serviceMediaType)
	if exitingServiceByMediaType != nil {
		err = &ErrServiceDuplicateIdentifier
		log.Fatal(s, "This mediatype '"+serviceMediaType+"' register already before for other service and not legal to reuse same mediatype for other services\n")
	} else {
		ss.poolByMediaType[serviceMediaType] = s
	}
	return
}

func (ss *Services) registerServiceByURI(s protocol.Service) (err protocol.Error) {
	var serviceURI = s.URI()
	if serviceURI != "" {
		var exitingServiceByURI, _ = ss.GetServiceByURI(serviceURI)
		if exitingServiceByURI != nil {
			err = &ErrServiceDuplicateIdentifier
			log.Fatal(s, "URI associated for '"+s.MediaType()+" service with `"+serviceURI+"` as URI, Used before for other service and not legal to reuse URI for other services\n"+
				"	Exiting service MediaType is: "+exitingServiceByURI.MediaType())
		} else {
			ss.poolByMediaType[serviceURI] = s
		}
	}
	return
}

// Services use to get all services registered.
//
//libgo:impl libgo/protocol.Services
func (ss *Services) Services() []protocol.Service { return ss.poolByRegisterTime }

// GetServiceByID use to get specific service handler by service ID
//
//libgo:impl libgo/protocol.Services
func (ss *Services) GetServiceByID(serviceID protocol.MediaTypeID) (ser protocol.Service, err protocol.Error) {
	ser = ss.poolByID[serviceID]
	if ser == nil {
		err = &ErrNotFound
	}
	return
}

// GetServiceByMediaType use to get specific service handler by service URI
//
//libgo:impl libgo/protocol.Services
func (ss *Services) GetServiceByMediaType(mt string) (ser protocol.Service, err protocol.Error) {
	ser = ss.poolByMediaType[mt]
	if ser == nil {
		err = &ErrNotFound
	}
	return
}

// GetServiceByURI use to get specific service handler by service URI path
//
//libgo:impl libgo/protocol.Services
func (ss *Services) GetServiceByURI(uri string) (ser protocol.Service, err protocol.Error) {
	ser = ss.poolByURIPath[uri]
	if ser == nil {
		err = &ErrNotFound
	}
	return
}

// DeleteService use to delete specific service in services list.
func (ss *Services) DeleteService(s protocol.Service) (err protocol.Error) {
	delete(ss.poolByID, s.ID())
	delete(ss.poolByMediaType, s.MediaType())
	delete(ss.poolByMediaType, s.URI())
	// TODO::: delete from ss.poolByRegisterTime
	return
}
