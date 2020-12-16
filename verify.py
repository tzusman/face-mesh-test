#!/bin/bash

# I'm python, I swear

SIZE=$(wc -c $2 | awk '{print $1}')

echo "{\"score\": ${SIZE}}"
