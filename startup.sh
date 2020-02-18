#!/bin/sh

APP_HOME=`dirname $0`
APP_HOME=`cd ${APP_HOME}; pwd`
APP_LAUNCHER="$APP_HOME/wechat-token"
APP_PID_FILE="$APP_HOME/app.pid"

APP_DAEMON_OUT="$APP_LOG_DIR/app.out"

echo -e "wechat-token starting ... "

if [ -f $APP_PID_FILE ]; then
  if kill -0 `cat $APP_PID_FILE` > /dev/null 2>&1; then
     echo wechat-token already running as process `cat $APP_PID_FILE`.
     exit 0
  fi
fi

nohup $APP_LAUNCHER > $APP_DAEMON_OUT 2>&1 < /dev/null &

if [ $? -eq 0 ]
then
  if /bin/echo -n $! > "$APP_PID_FILE"
  then
    sleep 1
    echo STARTED
  else
    echo FAILED TO WRITE PID
    exit 1
  fi
else
  echo SERVER DID NOT START
  exit 1
fi
