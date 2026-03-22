FROM ghcr.io/gleam-lang/gleam:v1.14.0-erlang-alpine AS dev

WORKDIR /app

COPY gleam.toml ./
RUN gleam deps download

COPY . .

EXPOSE 8082

# dev stage — used by devcontainer, keeps shell open
CMD ["sh"]

# ── production build ──
FROM dev AS build

RUN gleam export erlang-shipment

FROM erlang:27-alpine AS runtime

WORKDIR /app

COPY --from=build /app/build/erlang-shipment /app

EXPOSE 8082

CMD ["/app/entrypoint.sh", "run"]
