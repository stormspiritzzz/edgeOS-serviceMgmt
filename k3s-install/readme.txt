master������
./k3s-compose-master.sh

agent������
K3S_URL=https://myserver:6443 K3S_TOKEN=mynodetoken ./install.sh
���磺K3S_URL=https://192.168.2.109:6443 K3S_TOKEN=K10e3613c9e2748c9771ffebff3fb77231c8d9694c5c21af308f43903e85e55f18f::server:628fd3c7a3067d87cb8141ca21fe47ca ./k3s-compose-agent.sh

k3s�汾��
v1.18.6+k3s1

���Ի���docker�汾��
docker-ce | 5:19.03.2~3-0~ubuntu-xenial


�ȿӼ�¼��
1.master��agent��hostname����һ��������������agent����û���������е�����

2.k3sĬ��ʹ��containerd��Ϊ������agent������ָ��docker(INSTALL_K3S_EXEC="--docker")