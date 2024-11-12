FROM chainguard/wolfi-base:latest

# Add non-root user
USER nonroot

# Create necessary directories with correct permissions
WORKDIR /inbox451

# Copy binary
COPY --chown=nonroot:nonroot inbox451 .

# 8080: HTTP API
# 1025: SMTP Server
# 1143: IMAP Server
EXPOSE 8080 1025 1143

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Default environment variables
ENV INBOX451_SERVER_HTTP_PORT=":8080" \
    INBOX451_SERVER_SMTP_PORT=":1025" \
    INBOX451_SERVER_IMAP_PORT=":1143" \
    INBOX451_LOGGING_LEVEL="info" \
    INBOX451_LOGGING_FORMAT="json"

# Run the application
ENTRYPOINT ["./inbox451"]
