package server // import "github.com/DevanshMathur19/docker-v23/api/server"

import (
	"github.com/DevanshMathur19/docker-v23/api/server/httputils"
	"github.com/DevanshMathur19/docker-v23/api/server/middleware"
	"github.com/sirupsen/logrus"
)

// handlerWithGlobalMiddlewares wraps the handler function for a request with
// the server's global middlewares. The order of the middlewares is backwards,
// meaning that the first in the list will be evaluated last.
func (s *Server) handlerWithGlobalMiddlewares(handler httputils.APIFunc) httputils.APIFunc {
	next := handler

	for _, m := range s.middlewares {
		next = m.WrapHandler(next)
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		next = middleware.DebugRequestMiddleware(next)
	}

	return next
}
