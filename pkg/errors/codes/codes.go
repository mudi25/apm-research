package codes

type Code string

const (
	Success          Code = "SUCCESS"
	BadRequest       Code = "BAD_REQUEST"
	Unauthorized     Code = "UNAUTHORIZED"
	PermissionDenied Code = "PERMISSION_DENIED"
	DataNotFound     Code = "DATA_NOT_FOUND"
	Conflict         Code = "DATA_CONFLICT"
	PathNotFound     Code = "PATH_NOT_FOUND"
	MethodNotFound   Code = "METHOD_NOT_FOUND"
	Internal         Code = "INTERNAL_ERROR"
	Unavailable      Code = "UNAVAILABLE"
	UnknownError     Code = "UNKNOWN_ERROR"
)

func (c Code) HttpStatus() int {
	switch c {
	case Success:
		return 200
	case BadRequest:
		return 400
	case Unauthorized:
		return 401
	case PermissionDenied:
		return 403
	case DataNotFound, PathNotFound:
		return 404
	case MethodNotFound:
		return 405
	case Conflict:
		return 409
	case Unavailable:
		return 503
	default:
		return 500
	}
}
