FROM golang:latest

# 创建工作目录
RUN mkdir -p /data/www/ferry_ship/

# 进入工作目录
WORKDIR /data/www/ferry_ship

# 将当前目录下的所有文件复制到指定位置
COPY . /data/www/ferry_ship

# 下载beego和bee
ENV GIT_SSL_NO_VERIFY=1
RUN go get github.com/beego/bee/v2
# 端口
EXPOSE 8080
# 运行
CMD ["bee", "pack", "-be", "GOOS=linux", "-exs=.:node_modules:src:logs:conf/app.conf"]
