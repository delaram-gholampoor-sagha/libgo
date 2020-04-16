/* For license and copyright information please see LEGAL file in repository */

package services

import chaparkhane "../ChaparKhane"

// Init use to register all available services to given server.
func Init(s *chaparkhane.Server) {
	s.Services.RegisterService(&deleteIndexHashService)
	s.Services.RegisterService(&deleteIndexRecordService)
	s.Services.RegisterService(&deleteRecordService)
	s.Services.RegisterService(&findRecordsService)
	s.Services.RegisterService(&getIndexHashNumberService)
	s.Services.RegisterService(&getRecordService)
	s.Services.RegisterService(&listenToIndexService)
	s.Services.RegisterService(&readRecordService)
	// s.Services.RegisterService()
	s.Services.RegisterService(&setIndexService)
	s.Services.RegisterService(&setRecordService)
	s.Services.RegisterService(&warnAboutRecordService)
	s.Services.RegisterService(&writeRecordService)
	// s.Services.RegisterService()
}
