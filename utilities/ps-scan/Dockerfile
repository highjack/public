FROM debian:latest
RUN apt-get -y update
RUN apt-get -y install mono-complete wget zip
COPY ./regexes /regexes
RUN wget https://dist.nuget.org/win-x86-commandline/latest/nuget.exe -O /usr/bin/nuget.exe
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT /entrypoint.sh
