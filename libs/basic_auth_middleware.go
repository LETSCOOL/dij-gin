package libs

import (
	. "github.com/letscool/dij-gin"
	"net/http"
	"strconv"
)

const RefKeyForBasicAuthAccountCenter = "_.mdl.basic_auth_account"
const BasicAuthUserKey = "BasicAuthUserKey"

// AccountForBasicAuth
//
//	     ac := &AccountCenter{} // the struct implements AccountForBasicAuth interface
//		 config := NewWebConfig().SetDependentRef(RefKeyForBasicAuthAccountCenter, ac)
type AccountForBasicAuth interface {
	GetRealm() string
	SearchCredential(credential string) (account any, found bool)
}

type BasicAuthMiddleware struct {
	WebMiddleware

	account AccountForBasicAuth `di:"_.mdl.basic_auth_account"`

	realm string
}

func (b *BasicAuthMiddleware) DidDependencyInitialization() {
	realm := b.account.GetRealm()
	b.realm = "Basic realm=" + strconv.Quote(realm)
}

func (b *BasicAuthMiddleware) Authorize(ctx struct {
	WebContext `http:"basic_auth,method=handle"`
}) {
	authorization := ctx.GetRequestHeader("Authorization")
	user, found := b.account.SearchCredential(authorization)
	if !found {
		ctx.Header("WWW-Authenticate", b.realm)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set(BasicAuthUserKey, user)
}
