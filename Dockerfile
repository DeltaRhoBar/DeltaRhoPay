FROM golang:latest AS builder

WORKDIR /app

COPY . .

# RUN apt-get update && \
#     apt-get install -y\
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*

RUN go build -o ./tmp/delta_rho_pay ./cmd



FROM debian:latest

EXPOSE 8080

WORKDIR /app

COPY ./web ./web
COPY ./data ./data
COPY --from=builder /app/tmp/delta_rho_pay /usr/local/bin/

CMD ["/usr/local/bin/delta_rho_pay"]
