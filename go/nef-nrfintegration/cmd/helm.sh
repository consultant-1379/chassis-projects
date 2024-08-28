#!/usr/bin/env bash
helm del --purge nef-nrfintegration
helm lint ./helm/nef-nrfintegration
helm install ./helm/nef-nrfintegration -n nef-nrfintegration