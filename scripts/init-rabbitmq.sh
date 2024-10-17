#!/bin/bash
set -e

# Wait for RabbitMQ to start
until rabbitmqctl status; do
  echo "Waiting for RabbitMQ to start..."
  sleep 1
done

# Create user and set permissions
rabbitmqctl add_user $RABBITMQ_USER $RABBITMQ_PASSWORD
rabbitmqctl set_user_tags $RABBITMQ_USER administrator
rabbitmqctl set_permissions -p / $RABBITMQ_USER ".*" ".*" ".*"

echo "RabbitMQ user created and configured."