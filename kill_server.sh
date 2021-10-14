echo Trying to kill processes on port $PRODUCTION_PORT
sudo kill -9 $(sudo lsof -t -i:$PRODUCTION_PORT)