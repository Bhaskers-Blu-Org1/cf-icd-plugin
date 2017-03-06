FROM ubuntu
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install curl git -y
RUN curl https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz | tar -xz -C /usr/local
RUN mkdir /app
RUN mkdir /app/icd
COPY src /app/icd/src
WORKDIR /app/icd
ENV PATH /
ENV GOPATH /app/icd/vendor:/app/icd
RUN go
