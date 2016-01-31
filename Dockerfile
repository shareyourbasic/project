FROM quay.io/brianredbeard/corebox

ADD bin/app-linux-amd64 /opt/sercand/app

ADD public /opt/sercand/public/
ADD data.json /opt/sercand/data.json

WORKDIR /opt/sercand
EXPOSE 80
CMD ["/opt/sercand/app"]