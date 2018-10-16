FROM ai-image.jd.com/go/golang:1.10

WORKDIR /go/src
RUN mkdir -p DirMonitor
#挂载ceph
COPY .  DirMonitor/

RUN go install DirMonitor

WORKDIR DirMonitor

CMD go run *.go -c 汇天01 -m exporter -p /mnt/cephfs/algor-api/user/*/*,/mnt/cephfs/algor-api/dataset/user/*/*,/mnt/cephfs/algor-api/dataset/public/* -NF 10000 -NB 10488576
