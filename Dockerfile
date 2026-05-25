FROM golang:1.26-alpine AS base
RUN mkdir -p /empty

FROM base AS server_builder
WORKDIR /app
COPY ./server/go.mod ./server/go.sum ./
RUN go mod download

COPY ./server/ .
RUN go build -o makesweet-server .

FROM base AS makesweet_builder
RUN apk update && \
  apk add --no-cache \
  python3 \
  pipx \
  uv \
  binutils \
  scons \
  gcc \
  musl-dev \
  patchelf && \
  rm -rf /var/cache/apk/*

WORKDIR /app

COPY ./makesweet-py/pyproject.toml ./makesweet-py/uv.lock ./
RUN uv sync --frozen

COPY ./makesweet-py/src/ ./src/
COPY ./makesweet-py/makesweet-py.spec ./
RUN uv run pyinstaller makesweet-py.spec

RUN pipx run staticx dist/makesweet-py dist/makesweet-py-static --strip


FROM scratch
COPY --from=base /empty /tmp

WORKDIR /bin
COPY --from=server_builder /app/makesweet-server ./makesweet-server
COPY --from=makesweet_builder  --chmod=+x /app/dist/makesweet-py-static ./makesweet-py
COPY ./templates/ /templates/

EXPOSE 8080
ENTRYPOINT ["/bin/makesweet-server"]