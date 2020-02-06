docker_daemon_config = %q{'{"cluster-store":"consul://127.0.0.1:8500"}'}
docker = <<-EOF
count=$(apt list --installed | grep docker | wc -l)
if [ $count == 0 ]; then
  curl -fsSL https://get.docker.com/ | sh
  mkdir -p /etc/docker
  echo #{docker_daemon_config} > /etc/docker/daemon.json
  systemctl enable docker
  systemctl start docker
  usermod -a -G docker ubuntu
else
  echo docker is already installed
fi
EOF

consul_initial = <<-EOF
count=$(docker ps -q --filter name=consul | wc -l)
if [ $count == 0 ]; then
  docker run -d --name consul --network=host --restart always consul agent -server -bind 192.168.59.101 -bootstrap
fi
# wait for consul to come up
while ! curl -s localhost:8500 > /dev/null; do sleep 1; done
echo "Consul is available!"
EOF

consul_joiner = <<-EOF
count=$(docker ps -q --filter name=consul | wc -l)
if [ $count == 0 ]; then
 docker run -d --network=host --name consul --restart always consul agent -server -bind $IP -join 192.168.59.101
fi
# wait for consul to come up
while ! curl -s localhost:8500 > /dev/null; do sleep 1; done
echo "Consul is available!"
EOF
