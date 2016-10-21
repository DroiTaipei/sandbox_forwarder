#!/bin/bash

# prepare specific environment config
ln -s '/sandbox_forwarder/conf.d/'$PHASE'.toml' /sandbox_forwarder/conf.d/current.toml

# start sandbox_forwarder
/sandbox_forwarder/sandbox_forwarder -config /sandbox_forwarder/conf.d/current.toml