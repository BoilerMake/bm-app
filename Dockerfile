FROM golang:1.12 AS builder

WORKDIR /backend
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN make build

#RUN echo $ENV_MODE
#FROM scratch

#COPY --from=builder /backend/bin/server .

# Commented out until we have a web folder
#COPY --from=builder /backend/web/ .
ENV DB_CONN=user=postgres dbname=bm sslmode=disable

ENTRYPOINT /backend/bin/server
