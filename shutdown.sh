#!/bin/sh

APP_HOME=`dirname $0`
APP_HOME=`cd ${APP_HOME}; pwd`
APP_PID_FILE="$APP_HOME/app.pid"

echo -e "wechat-token stopping ... "

if [ ! -f "$APP_PID_FILE" ]
then
  echo "no app to stop (could not find file $APP_PID_FILE)"
else
  kill $(cat "$APP_PID_FILE")
  rm "$APP_PID_FILE"
  echo STOPPED
fi
exit 0
