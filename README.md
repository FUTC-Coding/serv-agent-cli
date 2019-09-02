# serv-agent-cli
## How to install

run this to build the docker image

```
sudo docker build \
--build-arg AGENT_USER='YOUR_DB_USERNAME' \
--build-arg AGENT_PASS='YOUR_DB_PASSWORD' \
--build-arg AGENT_IP='IP:PORT' \
-t serv-agent-cli .
```

Then run this to run the container and start the logging process

`sudo docker run --rm -it serv-agent-cli /serv-agent-cli start`
