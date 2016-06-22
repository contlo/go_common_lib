mkdir -p $HOME/goworkspace/log
if [ "$1"=="development" ]
then
  sed "s/var_go_env/$1/g" $GOPATH/src/go_common_lib/config/supervisord.conf | sed 's/user=deployer//g' > supervisord.conf
else
  sed "s/var_go_env/$1/g" $GOPATH/src/go_common_lib/config/supervisord.conf > supervisord.conf
fi
sudo supervisorctl reload
sudo supervisord -c supervisord.conf
