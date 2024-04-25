package web

import (
	"basic-go/week2/webook/internal/domain"
	"basic-go/week2/webook/internal/service"
	ijwt "basic-go/week2/webook/internal/web/jwt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	// 要求昵称长度在1到20个字符之间，禁止昵称为纯数字，禁止昵称为纯特殊符号或下划线
	nicknameRegexPattern = `^(?=.{1,20}$)(?!^[0-9]*$)(?!^[\\W_]*$)[a-zA-Z0-9\u4e00-\u9fa5\\._-]+$`
	birthdayRegexPattern = `^(19|20)\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`
	resumeRegexPattern   = `^.{1,200}$`

	bizLogin = "login"
)

// UserHandler 定义一个专门处理有关User的路由的Handler
type UserHandler struct {
	svc     service.UserService
	codeSvc service.CodeService
	ijwt.Handler

	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	nicknameRegexExp *regexp.Regexp
	birthdayRegexExp *regexp.Regexp
	resumeRegexExp   *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService,
	hdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		svc:     svc,
		codeSvc: codeSvc,
		Handler: hdl,

		// 预编译正则表达式提升性能
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		nicknameRegexExp: regexp.MustCompile(nicknameRegexPattern, regexp.None),
		birthdayRegexExp: regexp.MustCompile(birthdayRegexPattern, regexp.None),
		resumeRegexExp:   regexp.MustCompile(resumeRegexPattern, regexp.None),
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {

	// 使用Group分组路由，简化注册路由中的路径长度
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp) // note 此处是HandlerFunc类型，不需要写SignUp方法后的括号及参数（不然就是调用方法了）
	// ug.POST("/login", c.Login)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/logout", h.LogoutJWT)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
	ug.GET("/refresh_token", h.RefreshToken)
	ug.POST("/login_sms/code/send", h.SendSMSLog)
	ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SendSMSLog(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}

	// 调用service层的发送验证码
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "发送失败",
		})
	}
}
func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}
	// 验证码正确，调用service层进行登录
	// 因为用户可能未用手机号注册，所以需要调用FindOrCreate方法
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 登录成功后先设置refresh-token
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（JWT的refresh-token）")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})

}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	// note 习惯：使用 方法内部类 接收body的参数
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req Req
	// 调用Bind方法来反序列化 [1. 自动根据Content-Type进行绑定；2. 若有错误，自动返回到前端页面]
	// note Bind方法中一定要传req的地址！(源码有写)
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 校验
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入密码错误")
		return
	}
	isEmail, err := h.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（正则Pattern错误）")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}
	isPassword, err := h.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（正则Pattern错误）")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含至少8个字符，至少1个大写字母，1个小写字母，1个数字和1个特殊字符")
		return
	}

	// 调用 service层 【需要传入的对象是领域对象，而不是req】
	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	// note 处理邮箱相同的冲突err，即需要拿到 mysql 的唯一索引冲突
	// note 不能直接 if err==dao.ErrDuplicateEmail，因为web层里不能直接调dao层的东西，所以得一层层传，使得Handler只保持对service的依赖，避免跨层依赖
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱已被注册")
	default:
		ctx.String(http.StatusOK, "系统错误（web层的SignUp）")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		// 邮箱和密码
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
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
		return
	default:
		ctx.String(http.StatusOK, "系统错误（web层的Login）")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		// 邮箱和密码
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 登录成功后先设置refresh-token
		err = h.SetLoginToken(ctx, u.Id)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误（JWT的refresh-token）")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	default:
		ctx.String(http.StatusOK, "系统错误（web层的Login）")
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		Resume   string `json:"resume"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 对昵称、生日和个人简介进行正则规范
	isNickname, err := h.nicknameRegexExp.MatchString(req.Nickname)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（nickname的regex错误）")
		return
	}
	if !isNickname {
		ctx.String(http.StatusOK, "昵称不合法")
		return
	}
	isBirthday, err := h.birthdayRegexExp.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（birthday的regex错误）")
		return
	}
	if !isBirthday {
		ctx.String(http.StatusOK, "生日格式不合法")
		return
	}
	isResume, err := h.resumeRegexExp.MatchString(req.Resume)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（resume的regex错误）")
		return
	}
	if !isResume {
		ctx.String(http.StatusOK, "个人简介不合法")
		return
	}

	// 从session中取出userId【因为sess.get返回的interface类型，所以用类型断言认定为int64类型的值】
	// TODO 可以用 jwt？
	sess := sessions.Default(ctx)
	var userId int64
	if val, ok := sess.Get("userId").(int64); ok {
		userId = val
	} else {
		panic("userId的session的值不是int64")
	}

	// 除了用regex校验生日，还可以调用time.Parse方法【但返回的是time类型】
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式不合法")
		return
	}

	// 在web层调用service()，要用domain往下传
	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       userId,
		Nickname: req.Nickname,
		Birthday: birthday,
		Resume:   req.Resume,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误（web层的edit）")
		return
	}

	ctx.String(http.StatusOK, "更新成功（Edit）")

}

func (h *UserHandler) Profile(ctx *gin.Context) {

	// 方式一：从session中取出uid
	//sess := sessions.Default(ctx)
	//var userId int64
	//if val, ok := sess.Get("userId").(int64); ok {
	//	userId = val
	//} else {
	//	panic("userId的session的值不是int64")
	//}

	// 方式二：从uc中取
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	u, err := h.svc.FindById(ctx, uc.Uid)

	if err != nil {
		ctx.String(http.StatusOK, "系统错误（web层的profile）")
	}

	// 不能将domain.user直接传给前端，从中挑出nickname、Email、birthday和resume
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Birthday string `json:"birthday"`
		Resume   string `json:"resume"`
	}

	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		Birthday: u.Birthday.Format(time.DateOnly),
		Resume:   u.Resume,
	})

}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// 约定前端在 Authorization 里面带上 refresh_token
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RefreshKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 在redis中查看ssid
	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// redis有问题或者ssid存在（表明用户已退出） ==> redis有问题或者 token无效
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出成功",
	})
}
