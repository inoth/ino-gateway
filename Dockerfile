FROM alpine@sha256:e7d88de73db3d3fd9b2d63aa7f447a10fd0220b7cbf39803c803f2af9ba256b3
LABEL maintainer = "inoth" version = "v1.0" description = "ino-gateway"
EXPOSE 9000
EXPOSE 9001
WORKDIR /

RUN rm -f /etc/localtime \
&& ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache bash
RUN apk --no-cache add ca-certificates
RUN apk add dmidecode

COPY . .

RUN chmod +x /ino-gateway
ENTRYPOINT ["/ino-gateway"]
