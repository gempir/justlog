FROM alpine:3.6
COPY gempbotgo /
COPY configs /configs
CMD ["/gempbotgo"]
EXPOSE 8025