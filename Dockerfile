FROM golang:1.26-alpine AS base
RUN mkdir -p /empty

FROM base AS server_builder
WORKDIR /app
COPY ./server/go.mod ./server/go.sum ./
RUN go mod download

COPY ./server/ .
RUN go build -o makesweet-server .

FROM python:3.14-slim AS makesweet_builder
RUN apt-get update && \
  apt-get install -y \
  pipx \
  binutils &&\
  apt-get autoremove -y && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

ENV PATH="/root/.local/bin:${PATH}"

RUN pipx install uv

WORKDIR /app

COPY ./makesweet-py/pyproject.toml ./makesweet-py/uv.lock ./
RUN uv sync --frozen

COPY ./makesweet-py/src/ ./src/
COPY ./makesweet-py/makesweet-py.spec ./
RUN uv run pyinstaller makesweet-py.spec


FROM gcr.io/distroless/cc AS final
COPY --from=base /empty /tmp

WORKDIR /bin
COPY --from=server_builder /app/makesweet-server ./makesweet-server
COPY --from=makesweet_builder /app/dist/makesweet-py ./makesweet-py
COPY ./templates/ /templates/

EXPOSE 8080
ENTRYPOINT ["/bin/makesweet-server"]