#!/usr/bin/env bash

operator-sdk build vdksystem/db-operator && \

docker push vdksystem/db-operator:latest


