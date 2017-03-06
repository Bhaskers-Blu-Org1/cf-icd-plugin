FROM ubuntu
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install curl git gcc -y
RUN curl https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz | tar -xz -C /usr/local
RUN mkdir /app
RUN mkdir /app/icd
COPY src /app/icd/src
WORKDIR /app/icd
ENV PATH $PATH:/usr/local/go/bin 
ENV GOPATH /app/icd/vendor:/app/icd
RUN go get github.com/cloudfoundry/cli
ENV GOOS linux
ENV GOARCH amd64
RUN go build -o bin/linux64/icd icd
ENV GOOS darwin
ENV GOARCH amd64
RUN go build -o bin/osx64/icd icd
ENV GOOS windows
ENV GOARCH amd64
RUN go build -o bin/win64/icd icd
