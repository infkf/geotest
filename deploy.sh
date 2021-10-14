echo "Starting building the binary"
go build -o geotest_prod .
echo "Build finished"
# Killing server
ssh -i $KEYFILE_PATH $AWS_USER@$AWS_IP PRODUCTION_PORT=$PRODUCTION_PORT '. ~/projects/geotest/kill_server.sh'

# Copying binary and IP database files
echo "===== Copying files ======"
scp -i $KEYFILE_PATH -r ./db $AWS_USER@$AWS_IP:/home/$AWS_USER/projects/geotest/db
scp -i $KEYFILE_PATH ./geotest_prod $AWS_USER@$AWS_IP:/home/$AWS_USER/projects/geotest/geotest_prod
echo "===== Copying finished ======"
# Starting the server
ssh -i $KEYFILE_PATH $AWS_USER@$AWS_IP "cd ~/projects/geotest && (nohup ~/projects/geotest/geotest_prod 1>/dev/null 2>/dev/null &) ; echo Server started"

echo "Deploy finished"
