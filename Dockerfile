FROM centos:7

# install sidecar-injector binary
COPY build/_output/bin/adminssion-webhook /root/adminssion-webhook

# set entrypoint
CMD ["/root/adminssion-webhook"]
