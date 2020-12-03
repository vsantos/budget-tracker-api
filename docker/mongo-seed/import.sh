#! /bin/bash

mongoimport \
    --host mongodb \
    --db budget-tracker \
    --collection users \
    --type json \
    --file /mongo-seed/init.json \
    --jsonArray