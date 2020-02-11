FROM centos:6

RUN yum update -y && \
	yum groupinstall -y development && \
	yum install -y  \
	git \
	make \
	python2 \
	rpm-build \
	rpm-sign \
	wget \
	rng-tools \
	rpm-build \
	sudo \
	time \
	wget \
	yum-utils \
	tar \
	openssl \
	bzip2-devel \
	openssl-devel

ENV PATH="/root/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/go/bin"
ENV GOPATH=$HOME/go
RUN wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz -O - | tar -zx
RUN \
	mkdir -p $GOPATH/src/github.com/nutanix/clusters \
	mkdir -p $GOPATH/bin

RUN mkdir -p /usr/local/bin
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | INSTALL_DIRECTORY=/usr/local/bin sh

WORKDIR /home

COPY go.mod go.sum /home/
RUN go mod download

COPY . /home
WORKDIR /home

RUN go mod download

RUN \
       go build -o ./build/curate-clusters-service ./cmd/curate-clusters-service && \
       go build -o ./build/on-prem-agent ./cmd/on-prem-agent

CMD ["/home/build/curate-clusters-service"]
