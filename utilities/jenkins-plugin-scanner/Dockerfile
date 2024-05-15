FROM jenkins/jenkins:2.414.3-slim-jdk17
WORKDIR /downloads
RUN apt-get -y update
RUN apt-get install -y wget
RUN wget https://github.com/jenkins-zh/jenkins-cli/releases/latest/download/jcli-linux-amd64.tar.gz
RUN tar xzv jcli-linux-amd64.tar.gz
RUN mv jcli /usr/bin/
RUN wget https://github.com/aquasecurity/trivy/releases/download/v0.48.3/trivy_0.48.3_Linux-64bit.tar.gz
RUN tar xf trivy_0.48.3_Linux-64.tar.gz
RUN mv trivy /usr/bin
COPY get-token /usr/bin
ENTRYPOINT bash /download/entrypoint.sh


