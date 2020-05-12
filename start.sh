#!/usr/bin/env bash
export WHICH_DB_STRING="dbname=<your db name> host=<your db host> user=<your db username> password=<your password>" # add whatever you need to the connection string
export WHICH_GOOGLE_CLIENT_ID="<your google oauth client ID>"
export WHICH_GOOGLE_CLIENT_SECRET="<your google oauth client secret>"
export WHICH_GOOGLE_CALLBACK_URL="<your host>/auth/callback" # needs to be allowed in the Google Oauth Console
cd /home/username/which # absolute path to the top which directory
./which