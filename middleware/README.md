## 基础库

### middleware [中间件]

- [x] comm - 通用中间件
- [ ] gin - 中间件
    - [x] cache - 缓存
    - [x] commmonlog - 日志
        - [x] ginlog - 请求日志
        - [x] rotate - 请求日志，自动切分
    - [x] gin-nice-recovery - 优雅恢复
    - [x] request_id - X-Requset—Id请求ID，方便日志trace
    - [x] jwt - TokenAuth
    - [x] expvar - expvar监控
    - [x] secure - 访问安全
    - [x] sentry - sentry监控
    - [x] cors - 跨域访问控制（Cross-Origin Resource Sharing）
    - [x] csrf - 跨站请求伪造（Cross-site request forgery）
    - [x] revision - 版本号（Response Headers: X-Revision-Id）
    - [x] location - 暴露服务器的hostname and scheme
    - [x] limit - 访问请求并发限制（接口粒度）
    - [x] limit-by-key - 访问请求并发限制（自定义粒度）
    - [x] redis-ip-limiter - 基于Redis的全局访问请求并发限制
    - [x] access-limit - 请求来源限制
    - [x] stats - API请求统计
    - [x] gindump - Debug模式下日志输出