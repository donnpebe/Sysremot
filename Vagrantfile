# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "centos65-x86_64"
  config.vm.box_url = "https://github.com/2creatives/vagrant-centos/releases/download/v6.5.1/centos65-x86_64-20131205.box"
  config.vm.provision "shell", inline: "mkdir -p /home/vagrant/go"
  config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/donnpebe/sysremot"
  config.vm.provision "shell", inline: "chown -R vagrant:vagrant /home/vagrant/go"
  install_go = <<-BASH
if [ ! -d "/usr/local/go" ]; then cd /tmp && wget https://storage.googleapis.com/golang/go1.4.1.linux-amd64.tar.gz && cd /usr/local && tar xvzf /tmp/go1.4.1.linux-amd64.tar.gz && echo 'export GOPATH=/home/vagrant/go; export PATH=/usr/local/go/bin:$PATH:$GOPATH/bin' >> /home/vagrant/.bashrc && su - vagrant -c 'go get github.com/tools/godep && cd $GOPATH/src/github.com/donnpebe/sysremot && godep restore && go build' && mv /home/vagrant/go/src/github.com/donnpebe/sysremot/sysremot /usr/local/bin/sysremot && chmod u+x /usr/local/bin/sysremot && service redis start; fi;
BASH
  config.vm.provision "shell", inline: 'yum install -y redis wget'
  config.vm.provision "shell", inline: install_go
end