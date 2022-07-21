FROM debian

WORKDIR /

COPY qinchy-scheduler /usr/local/bin

CMD ["/usr/local/bin/qinchy-scheduler"]