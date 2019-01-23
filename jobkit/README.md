jobkit
======

This package is meant to be a suite of helpers to make writing robust job workers easier.

It provides facilities that plug into `go-sdk/cron` to help with:
- A management server to streamline allowing forced runs of jobs.
- Sending email notifications for job results.
- Sending slack notifications for job results.
