# lurra

## Before deploying

- **Log the client IP.** Add the chi `ClientIPFrom*` middleware that matches your deployment
  (see [Client IP](https://go-chi.io/#/pages/client_ip)), before the request logger:
  - exposed directly: `middleware.ClientIPFromRemoteAddr`
  - behind nginx, Cloudflare, etc.: `middleware.ClientIPFromHeader("X-Real-IP", <proxy CIDRs>)`
  - behind N load balancers: `middleware.ClientIPFromXFFTrustedProxies(n)`

  httplog's concise schema drops the IP, so also add a middleware after the logger that calls
  `httplog.SetAttrs(ctx, slog.String("client.ip", middleware.GetClientIP(ctx)))`.
