FROM ubuntu:latest

RUN apt update && apt install -y curl git && curl -L https://foundry.paradigm.xyz | bash && ~/.foundry/bin/foundryup
ENV PATH="/root/.foundry/bin:${PATH}"

COPY anvile /usr/local/bin/anvile
RUN chmod +x /usr/local/bin/anvile
ENV PORT=80

ENTRYPOINT ["/usr/local/bin/anvile"]