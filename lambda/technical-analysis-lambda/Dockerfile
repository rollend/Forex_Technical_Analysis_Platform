FROM golang:1.21.0-bullseye as build-image

WORKDIR /go/src
COPY go.mod main.go ./

# Download all the dependencies
RUN go mod download

COPY . .

RUN go build -o ../bin

FROM public.ecr.aws/lambda/go:1

COPY --from=build-image /go/bin/ /var/task/

# Command can be overwritten by providing a different command in the template directly.
CMD ["technical-analysis-lambda"]
