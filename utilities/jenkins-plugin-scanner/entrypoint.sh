#!/bin/bash
get-token
rm -r /download/plugins 2>/dev/null
mkdir /download/plugins 2>/dev/null
for plugin in $(cat /download/jenkins_plugins.txt);do
jcli plugin install $plugin;
done
