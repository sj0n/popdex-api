package pokemon

import (
	"encore.dev/beta/errs"
	"encore.dev/middleware"
)

//encore:middleware target=all
func ETagMiddleware(req middleware.Request, next middleware.Next) middleware.Response {

	ifNoneMatch := req.Data().Headers["If-None-Match"]
	resp := next(req)
	var etag string
	
	switch payload := resp.Payload.(type) {
    case *PokemonProfile:
		etag = payload.ETag
    case *GroupByVersion:
		etag = payload.ETag
    case *PokemonLocations:
        etag = payload.ETag
    }

	if len(ifNoneMatch) != 0 && ifNoneMatch[0] == etag {
		return middleware.Response{HTTPStatus: 304, Err: errs.B().Msg("Not Modified").Err()}
	}

	return resp
}