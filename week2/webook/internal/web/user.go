package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 定义一个专门处理有关User的路由的Handler
type UserHandler struct {
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

func NewUserHandler() *UserHandler {
	return &UserHandler{
		// 预编译正则表达式提升性能
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

// RegisterRoutes 用于注册路由
func (c *UserHandler) RegisterRoutes(server *gin.Engine) {

	// 使用Group分组路由，简化注册路由中的路径长度
	ug := server.Group("/users")
	ug.POST("/signup", c.SignUp) // TODO 注意此处HandlerFunc类型，不需要写SignUp后的括号及参数
	ug.POST("/login", c.Login)
	ug.POST("/edit", c.Edit)
	ug.POST("/profile", c.Profile)

}

// SignUp 定义UserHandler上的方法作为应路由的的处理逻辑
func (c *UserHandler) SignUp(ctx *gin.Context) {
	// 习惯：优先使用方法内部类
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	// 调用Bind方法 [1. 自动根据Content-Type进行绑定；2. 若有错误，自动返回到前端页面]
	// TODO Bind方法中一定要传req的地址！！！！不然绑定不成功
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 校验
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入密码错误")
		return
	}

	isEmail, err := c.emailRegexExp.MatchString(req.Email)

	if err != nil {
		ctx.String(http.StatusOK, "系统错误（正则Pattern错误）")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}
	isPassword, err := c.passwordRegexExp.MatchString(req.Password)

	if err != nil {
		ctx.String(http.StatusOK, "系统错误（正则Pattern错误）")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含至少8个字符，至少1个大写字母，1个小写字母，1个数字和1个特殊字符")
		return
	}

	ctx.String(http.StatusOK, "SignUp的响应")
}

func (c *UserHandler) Login(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Login的响应")

}

func (c *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Edit的响应")

}

func (c *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Profile的响应")

}
