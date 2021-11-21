FROM scratch
COPY kollect /

ENTRYPOINT ["/kollect"]
