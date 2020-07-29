package app

// success response
// swagger:response success
type success struct {
	// in:body
	Body event
}


// Error Bad Request
// swagger:response badRequest
type swaggErrBadReq struct {
	// in:body
	Body struct {
		// HTTP status code 400 - Status Bad Request
		Code int `json:"code"`
		// Detailed error message
		Message string `json:"message"`
	}
}

// Error Forbidden
// swagger:response forbidden
type swaggErrForbidden struct {
	// in:body
	Body struct {
		// HTTP status code 403 - Forbidden
		Code int `json:"code"`
		// Detailed error message
		Message string `json:"message"`
	}
}

// Error Not Found
// swagger:response notFound
type swaggErrNotFound struct {
	// in:body
	Body struct {
		// HTTP status code 404 - Not Found
		Code int `json:"code"`
		// Detailed error message
		Message string `json:"message"`
	}
}

// Error Conflict
// swagger:response conflict
type swaggErrConflict struct {
	// in:body
	Body struct {
		// HTTP status code 409 - Conflict
		Code int `json:"code"`
		// Detailed error message
		Message string `json:"message"`
	}
}

// Error Interval Server
// swagger:response internal
type swaggErrInternal struct {
	// in:body
	Body struct {
		// HTTP status code 500 - Internal server error
		Code int `json:"code"`
		// Detailed error message
		Message string `json:"message"`
	}
}


