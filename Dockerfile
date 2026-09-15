FROM debian:stable-slim

# 时区/日志时间使用本地时间并携带 CA 证书（HTTPS 根证书）
ENV TZ=Asia/Shanghai \
    LANG=C.UTF-8

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

# 复用预编译产物（x86-64 动态链接 ELF，依赖 glibc，debian:stable-slim 满足其 GLIBC_2.34 上限）
COPY output/disrecall-linux-amd64 /app/disrecall
RUN chmod +x /app/disrecall

# 配置目录(config.yml/sqlite.db/logs)与下载文件目录由 docker-compose 挂载持久化
VOLUME ["/app/config", "/app/files"]

ENTRYPOINT ["./disrecall"]