# Step 1
FROM alpine:3.19.1
WORKDIR /opt/ats
ARG ATS_VERSION
ARG ATS_REVISION
ARG GO_VERSION="go1.22.0"
ARG LD_PATH="ats/src/config"
ARG LD_FLAGS="-X $LD_PATH.Version=${ATS_VERSION} -X $LD_PATH.GoVersion=${GO_VERSION} -X $LD_PATH.GitCommit=${ATS_REVISION}"
RUN apk add --no-cache --virtual .build-depsgcc libc-dev g++ make

RUN wget https://dl.google.com/go/${GO_VERSION}.linux-amd64.tar.gz
RUN tar xf ${GO_VERSION}.linux-amd64.tar.gz -C /usr/local
ENV PATH=$PATH:/usr/local/go/bin
RUN go version

# Install upx 
RUN wget https://github.com/upx/upx/releases/download/v4.2.3/upx-4.2.3-amd64_linux.tar.xz
RUN tar xf upx-4.2.3-amd64_linux.tar.xz -C /usr/local
ENV PATH=$PATH:/usr/local/upx-4.2.3-amd64_linux
RUN upx -h

COPY . /opt/ats

# -trimpath 移除源代码中的文件路径信息
# -ldflags -s：不生成符号表 -w：不生成DWARF调试信息
RUN go build -trimpath  -ldflags "-s -w $LD_FLAGS" -o output/ats src/main/main.go
RUN upx -9 output/ats
RUN output/ats -version
RUN cp ats.yaml output
RUN chmod +x docker-entrypoint.sh
RUN cp docker-entrypoint.sh output/entrypoint.sh

# Step 2
FROM alpine:3.19.1
ARG ATS_REVISION
ARG ATS_VERSION

COPY --from=0 /opt/ats/output/* /usr/local/bin

ENV ATS_REVISION=$ATS_REVISION \
    ATS_VERSION=$ATS_VERSION

ENTRYPOINT ["entrypoint.sh"]
