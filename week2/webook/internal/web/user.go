package web

import (
	"basic-go/week2/webook/internal/domain"
	"basic-go/week2/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// UserHandler 定义一个专门处理有关User的路由的Handler
type UserHandler struct {
	svc *service.UserService

	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
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

	// 调用 service层 【需要传入的对象是领域对象，而不是req】
	err = c.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	// TODO 处理邮箱相同的冲突err，即需要拿到 mysql 的唯一索引冲突
	// 不能直接 if err!=dao.ErrDuplicateEmail，因为web层里不能直接调dao层的东西，所以得一层层传
	// 使得Handler之保持对service的依赖，避免跨层依赖
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱已被注册")
	default:
		ctx.String(http.StatusOK, "系统错误（web层的SignUp）")
	}
}

func (c *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		// 邮箱和密码
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := c.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:
		// 登录成功后获取session，存入域对象u的id，便于profile和edit方法获取
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 15min
			MaxAge: 900,
		})
		err := sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误（保存登录状态的session）")
		}

		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "账号或密码错误")
	default:
		ctx.String(http.StatusOK, "系统错误（web层的Login）")
	}
}

func (c *UserHandler) Edit(ctx *gin.Context) {
	var profile domain.Profile
	if err := ctx.Bind(&profile); err != nil {

	}
	ctx.String(http.StatusOK, "Edit的响应")

}

func (c *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Profile的响应")

}
