#!/bin/bash
echo $(PWD)/volume
mkdir -p $(PWD)/volume
docker run -p 27017:27017 -v $(PWD)/volume:/data/db mongo
