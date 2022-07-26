/* For license and copyright information please see LEGAL file in repository */

package protocol

// Errors is the interface that must implement by any Application.
type Errors interface {
	RegisterError(err Error)
	GetErrorByID(id MediaTypeID) (err Error)
	GetErrorByMediaType(mt string) (err Error)
}

// Error is similar to opaque error model describe here: https://dave.cheney.net/paste/gocon-spring-2016.pdf
// or this RFC: https://tools.ietf.org/html/rfc7807
type Error interface {
	// Init() must call protocol.App.RegisterError() to register the error in application
	// Init(mediatype string)

	// Check both flat or chain situation.
	Equal(Error) bool

	// Who cause the error?
	// Internal	: means calling process logic has runtime bugs like HTTP server error status codes ( 500 – 599 ).
	// Caller	: Opposite of internal that indicate caller give some data that cause the error like HTTP client error status codes ( 400 – 499 )
	Internal() bool
	Temporary() bool // opposite is permanent situation

	// Notify error to user by graphic, sound and vibration (Haptic Feedback)
	Notify()

	MediaType
	Stringer
}
