FROM gcr.io/distroless/static

ARG TARGETPLATFORM

ADD bin/${TARGETPLATFORM} /bin
