FROM golang:1.12
WORKDIR $GOPATH/src/serv-agent-cli
ADD . /go/src/serv-agent-cli
RUN go get github.com/mackerelio/go-osstat/network
RUN go get github.com/shirou/gopsutil/cpu
RUN go get github.com/shirou/gopsutil/disk
RUN go get github.com/shirou/gopsutil/mem
RUN github.com/spf13/cobra
RUN go install /go/src/serv-agent-cli
ENTRYPOINT /go/bin/serv-agent-cli
