#!/bin/bash

if [ ! -f vendor ]; then
  glide up -u -s
fi

go build .
