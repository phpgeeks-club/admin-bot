# Run
# docker build -f ./Dockerfile -t geeksonator:latest .
# docker run -d --env-file=/path/to/.env --name geeksonator.app geeksonator:latest .

##################################
# STEP 1 build executable binary #
##################################

FROM golang:1.23.6-alpine as builder

LABEL org.opencontainers.image.source="https://github.com/phpgeeks-club/admin-bot"

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata make && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app

COPY . .

# Build the binary.
RUN make build

##############################
# STEP 2 build a small image #
##############################

FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /app/bin/geeksonator /app/geeksonator

# Use an unprivileged user.
USER appuser:appuser

# Run the binary.
ENTRYPOINT ["/app/geeksonator"]
