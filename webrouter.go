// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import "github.com/gin-gonic/gin"

type WebRouter interface {
	gin.IRouter
	BasePath() string
}

type WebRoutes interface {
	gin.IRoutes
	BasePath() string
}
