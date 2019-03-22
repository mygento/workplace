#!/usr/bin/env bash
export USERID=$(id -u)

echo -e "\033[37;1;42m \033[5mCleaning app/code \033[0m"
rm -fR ./application/app/code
echo
echo -e "\033[37;1;42m \033[5mCleaning app/design \033[0m"
echo
rm -fR ./application/app/design

docker-compose up -d
if [ -f ./application/bin/magento ]; then
  chmod +x ./application/bin/magento
fi

gulp
