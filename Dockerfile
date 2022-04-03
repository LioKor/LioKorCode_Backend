FROM golang:alpine

COPY . /project

WORKDIR /project

RUN go build -o build/go_main_service ./cmd/main.go
RUN ls | grep -v 'go' | grep -v 'Docker' | grep -v 'Readme' | grep -v 'build' | grep -v 'public' | xargs rm -rfv
RUN ls | grep -v 'build' | grep -v 'public' | xargs rm

CMD ["/project/build/go_main_service"]