Vagrant.configure(2) do |config|
  [
    ['192.168.0.11', 'master', '7001'],
    ['192.168.0.12', 'node1', '7002'],
    ['192.168.0.13', 'node2', '7003']
  ].collect.each_with_index do |data|
    config.vm.define data[1] do |node|
      node.vm.hostname = data[1]
      node.vm.box = 'bento/ubuntu-16.04'
      node.vm.network :private_network, ip: data[0]
      node.vm.network 'forwarded_port', guest: 7000, host: data[2]
      if data[1] == 'master'
        node.vm.network 'forwarded_port', guest: 80, host: 8080
      end
      node.vm.provider 'virtualbox' do |vb|
        vb.memory = 1024
        vb.cpus = 1
      end
      node.vm.provision 'shell', path: 'contrib/vagrant/init.sh'
    end
  end
end
