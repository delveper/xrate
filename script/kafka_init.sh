#!/bin/bash

CLUSTER_ID=$(/bin/kafka-storage random-uuid)
export CLUSTER_ID
exec /etc/confluent/docker/run