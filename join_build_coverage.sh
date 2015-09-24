#!/bin/bash
cat /home/ubuntu/*_coverage.out | grep -v "mode: set" >> /home/ubuntu/coverage.out
