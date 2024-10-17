FROM rabbitmq:3-management-alpine
COPY scripts/init-rabbitmq.sh /init-rabbitmq.sh
CMD ["/bin/bash", "-c", "/init-rabbitmq.sh && rabbitmq-server"]