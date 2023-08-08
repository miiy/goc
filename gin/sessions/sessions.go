package sessions

//
//import (
//	"github.com/gin-contrib/sessions"
//	"github.com/gin-contrib/sessions/redis"
//	"github.com/gin-gonic/gin"
//)
//
//func Sessions(name string) gin.HandlerFunc {
//	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
//	return sessions.Sessions(name, store)
//}
//
//func Default(c *gin.Context) sessions.Session {
//	return sessions.Default(c)
//}
