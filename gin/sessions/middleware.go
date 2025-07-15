package sessions

//
//import (
//	"github.com/gin-gonic/gin"
//	"log"
//	"net/http"
//)
//
//func SessionAuthenticationMiddleware(session *session.Session) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		sess, err := session.Store.Get(c.Request, session.Options.SessionCookie)
//		if err != nil {
//			log.Println(err)
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//
//		user := sess.Values["auth"]
//		if user == nil {
//			c.Redirect(http.StatusMovedPermanently, "/signin")
//			return
//		}
//		c.Set("auth", user)
//		c.Next()
//
//	}
//}
//
//func SessionFlashMiddleware(session *session.Session) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		sess, err := session.Store.Get(c.Request, session.Options.SessionCookie)
//		if err != nil {
//			log.Println(err)
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//
//		flashes := sess.Flashes()
//		if len(flashes) > 0 {
//			for _, v := range flashes {
//				if e, ok := v.(pkgValidator.ValidationErrorsTranslations); ok {
//					c.Set("flashes", e)
//				}
//			}
//
//			if err = sess.Save(c.Request, c.Writer); err != nil {
//				log.Println(err)
//				c.AbortWithStatus(http.StatusInternalServerError)
//				return
//			}
//		}
//
//		c.Next()
//
//	}
//}
//
