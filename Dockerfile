FROM golang:1.12
WORKDIR $GOPATH/src/serv-agent-cli
ADD . /go/src/serv-agent-cli
RUN go get github.com/mackerelio/go-osstat/network
RUN go get github.com/shirou/gopsutil/cpu
RUN go get github.com/shirou/gopsutil/disk
RUN go get github.com/shirou/gopsutil/mem
RUN go get github.com/spf13/cobra/cobra
RUN go get github.com/go-sql-driver/mysql
RUN go install /go/src/serv-agent-cli
ARG AGENT_USER
ARG AGENT_PASS
ARG AGENT_IP
RUN /go/bin/serv-agent-cli database -u $AGENT_USER -p $AGENT_PASS -i $AGENT_IP
ENTRYPOINT /go/bin/serv-agent-cli $1
