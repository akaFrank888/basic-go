## Golang学习（ready for 2026校招）

#### 启动：

1） docker compose：在webook目录下，执行`docker compose up`

2） 后端

3） 前端启动： 在webook-fe目录下，执行`npm install`，等下载好接着执行`npm run dev`

4） 若增删了模块的方法，可能要在webook目录下，执行`wire`

#### 代办：

1） 做一个该课程的最小可运行的工具包，当作笔记或者自己的代码库，以方便开发和总结



### Week 1 Go基础语法
学习时间： 2024/4/3 —— 2024/4/6

学习内容：（Go基本语法总结）
<img src="images/week1_syntax.png">

待提升：
1. 搭建自己的 <b>Go泛型工具库<b>
<img src="images/week1_generic_kit.png">

学习参考：https://github.com/ecodeclub/ekit/tree/dev/internal/slice

### Week 2 Gin、GORM入门与用户注册登录功能实现

学习时间： 2024/4/7 —— 2024/4/10

学习内容：（Gin和Gorm）
1. 如图
   <img src="images/week2_content.png">

2. 完成了handler层中edit()和profile()的编写
3. 基于cookie的实现保存session数据

作业截图：【接口测试工具：Apifox】

1. edit接口测试
   <img src="images/week2_edit.png">

2. profile接口测试
   <img src="images/week2_profile.png">

### Week 3 JWT、Redis入门和Kubernetes部署实战

学习时间： 2024/4/10 —— 2024/4/14

学习内容：
1. 如图
   <img src="images/week3_content.png">

2. <font color='brown'>跳过了“用k8s部署web、mysql和redis”，进度在“k8s部署web服务器”视频的29:33处。</font>
【跳过原因：1. k8s命令琐碎，跳过不影响后续课程 2. 面试时重点在回答概念，且入职后确保自己在公司接触到k8s的时候，能够看懂公司的k8s部署配置就行】

<img src="images/week3_content_k8s.png">

### Week 4 接口抽象技巧和短信服务

学习时间： 2024/4/15 —— 2024/4/19

学习内容：
1. 上图的“压测与缓存机制”
2. 短信服务实现
   <img src="images/week4_message.png">
3. <font color='brown'>第4周作业未提交 【涉及到普通锁】</font>
4. 用wire改造代码（依赖注入） + 面向接口编程

<img src="images/week4_wire.png">

### Week 5 单元测试和集成测试、第三方服务调用治理

学习时间： 2024/4/20 —— 2024/4/21

学习内容：
1. 单元测试和集成测试、第三方服务调用治理

学习进度：
1. <font color='brown'>跳过了Handler(下)及后面的内容</font>

### Week 6 OAuth2与微信扫码登录实现

学习时间： 2024/4/21 —— 2024/4/

学习内容：
1. 微信扫码登录

<img src="images/week6_wxlogin_process.png">

2. 长短token的实现
3. viper:配置模块
4. 日志模块

学习进度：
1. 微信登录所需要的APP_ID和APP_SECRET暂无申请到，可以考虑用别的 Oauth2 的，比如说支付宝，github，google，基本思路都一样，就是调用的 API 不同【待用其他平台的oauth2重构】
2. 

待更。。。。