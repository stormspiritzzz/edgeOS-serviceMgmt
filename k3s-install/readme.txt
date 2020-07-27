master启动：
./k3s-compose-master.sh

agent启动：
K3S_URL=https://myserver:6443 K3S_TOKEN=mynodetoken ./install.sh
例如：K3S_URL=https://192.168.2.109:6443 K3S_TOKEN=K10e3613c9e2748c9771ffebff3fb77231c8d9694c5c21af308f43903e85e55f18f::server:628fd3c7a3067d87cb8141ca21fe47ca ./k3s-compose-agent.sh

k3s版本：
v1.18.6+k3s1

测试环境docker版本：
docker-ce | 5:19.03.2~3-0~ubuntu-xenial


踩坑记录：
1.master和agent的hostname不能一样，不报错，但是agent机子没有正常运行的容器

2.k3s默认使用containerd作为容器，agent启动须指定docker(INSTALL_K3S_EXEC="--docker")