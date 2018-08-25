FROM alpine
COPY ./xds-api-test /root/
RUN chmod +x /root/xds-api-test
ENTRYPOINT ["/root/xds-api-test"]

