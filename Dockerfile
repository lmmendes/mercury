FROM chainguard/wolfi-base:latest

# Add non-root user
USER nonroot

# Create necessary directories with correct permissions
WORKDIR /inbox451

# Copy binary and config
COPY --chown=nonroot:nonroot inbox451 .
COPY --chown=nonroot:nonroot config.yml .

# 8080: HTTP API
# 1025: SMTP Server
# 1143: IMAP Server
EXPOSE 8080 1025 1143

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Environment variables
ENV INBOX451_SERVER_HTTP_PORT=":8080" \
  INBOX451_SERVER_SMTP_PORT=":1025" \
  INBOX451_SERVER_IMAP_PORT=":1143"

# Run the application
ENTRYPOINT ["./inbox451"]
CMD ["--config", "config.yml"]
