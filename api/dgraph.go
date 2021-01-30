package api

import "github.com/gin-gonic/gin"

//要不要做权限验证？ 权限验证实际是在dgraph后端，但是权限相关的错误码要在这里做映射。
func DgraphProxy(c *gin.Context) {

}
