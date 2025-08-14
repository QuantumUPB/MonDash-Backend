#!/bin/sh

/app/mondash &
exec nginx -g 'daemon off;'

