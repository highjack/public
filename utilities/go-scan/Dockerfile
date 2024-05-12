FROM golang:1.19-alpine
RUN apk --no-cache add ca-certificates git
ENV PYTHONUNBUFFERED=1
RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
RUN apk add wget 
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools
RUN apk add zip
WORKDIR /app
COPY ./app/run.py /app
RUN wget https://github.com/sonatype-nexus-community/nancy/releases/download/v1.0.42/nancy-v1.0.42-linux-amd64 -O nancy
RUN chmod +x ./nancy
CMD ["python3", "run.py"]
