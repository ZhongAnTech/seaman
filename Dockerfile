FROM golang:1.12-alpine3.10 as base
RUN apk --update upgrade
RUN apk --no-cache add tzdata bash curl busybox-extras make g++ libstdc++ git
RUN rm -rf /var/cache/apk/*
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.12.10/bin/linux/amd64/kubectl && chmod +x ./kubectl && mv ./kubectl /usr/local/bin/kubectl
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone
ENV GOPROXY https://goproxy.io


FROM base as builder
WORKDIR /tmp/seaman
COPY . .
RUN CGO_ENABLED=1 INSTALL_DIR=/seaman make install clean


FROM builder
WORKDIR /seaman
COPY --from=builder /seaman .
EXPOSE 8080
CMD ["./seaman"]