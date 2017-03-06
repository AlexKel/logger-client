FROM golang:1.8-alpine

RUN mkdir -p /go/src/gitlab.devops.ukfast.co.uk/fastbase/logger-client
COPY . /go/src/gitlab.devops.ukfast.co.uk/fastbase/logger-client
WORKDIR /go/src/gitlab.devops.ukfast.co.uk/fastbase/logger-client

RUN apk add --no-cache git
RUN go-wrapper download
RUN go-wrapper install

# Add goblin BDD framework
RUN mkdir -p /go/src/github.com/franela/goblin
RUN git clone https://github.com/franela/goblin /go/src/github.com/franela/goblin

# Add gomega BDD assertions
RUN mkdir -p /go/src/github.com/onsi/gomega
RUN git clone https://github.com/onsi/gomega /go/src/github.com/onsi/gomega

# Add gocov for test coverage
RUN mkdir -p /go/src/github.com/axw/gocov
RUN git clone https://github.com/axw/gocov /go/src/github.com/axw/gocov

RUN go get gopkg.in/yaml.v2
