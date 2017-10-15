#!/bin/bash

export DEBIAN_FRONTEND=noninteractive
HOSTNAME=$(hostname)

apt-get update > /dev/null
apt-get install nginx -y
rm /etc/nginx/sites-enabled/default

if [ $HOSTNAME = "master" ]; then
	cp /vagrant/contrib/vagrant/datalock-master.conf /etc/nginx/conf.d/
else
	cp /vagrant/contrib/vagrant/datalock-node.conf /etc/nginx/conf.d/
	mkdir -p /opt/datalock/ /opt/datalock/database
	chmod 0777 /opt/datalock/database
	cp /vagrant/bin/datalock /opt/datalock/
	cp -r /vagrant/public /opt/datalock/
	cp -r /vagrant/templates /opt/datalock/
	cp /vagrant/contrib/init/systemd/datalock.service /etc/systemd/system/
	systemctl enable datalock
	systemctl start datalock
fi

nginx -s reload
