/* For license and copyright information please see LEGAL file in repository */

package ss

import chaparkhane  "../ChaparKhane"

var transferConnectionService = chaparkhane.Service{
	Name:       "TransferConnection",
	IssueDate:  0,
	ExpiryDate: 0,
	Status:     chaparkhane.ServiceStatePreAlpha,
	Handler:    TransferConnection,
	Description: []string{
		`Use to transfer exiting connection from other server if related service exist in platform!
		 Connection can be transfer if server have free in connection pool!`,
	},
	TAGS: []string{},
}

// TransferConnection Use to transfer exiting connection from other server if related service exist in platform!
// Connection can be transfer if server have free in connection pool!
func TransferConnection(s *chaparkhane.Server, st *chaparkhane.Stream) {
	// If it is not register guest connection service, so

	// Check user can open new connection in this server first!

	// get Connection detail from s.Manifest.TechnicalInfo.AuthorizationServer
	// Call related SDK
}

type transferConnectionReq struct {
}
