# gosible

A wrapper around ansible to enable additional support for cool stuff.

loosely based off https://github.com/blueboxgroup/ursula-cli

```bash
$ go build

$ ./gosible help playbook run

Gosible playbook is a wrapper around ansible-playbook that adds some
additional useful features.

Usage:
  gosible playbook run [flags] file.yml [ansible-playbook arguments]

Flags:
  -e, --environment string        directory that contains ansible inventory
  -h, --help                      help for run
      --known-hosts-file string   location of known hosts file
  -s, --ssh-config-file string    Path to ssh config file to use.
  -f, --ssh-forward-agent         path to ssh config file to use

Global Flags:
      --config string   config file (default is $HOME/.gosible.yaml)



$ ./gosible playbook run -e tests/functional/environment \
    tests/functional/playbook/ping.yml  --become  

running: ansible_playbook --inventory tests/functional/environment/hosts tests/functional/playbook/ping.yml --become
PLAY [ensure connectivity to all nodes] ****************************************
TASK [Gathering Facts] *********************************************************
fatal: [127.0.0.1]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: ssh: connect to host 127.0.0.1 port 22: Connection refused\r\n", "unreachable": true}
	to retry, use: --limit @/home/pczarkowski/development/go/src/github.com/paulczar/gosible/tests/functional/playbook/ping.retry
PLAY RECAP *********************************************************************
127.0.0.1                  : ok=0    changed=0    unreachable=1    failed=0   
There was an error running Ansible:  exit status 4
```

> Note: the above as written is expected to fail the actual ansible run.

If you have Vagrant installed you can exercise the vagrant provisioner (Vagrantfile needs to be in current working dir, or in the environment dir):

```
$ ./gosible playbook run --provisioner=vagrant -e tests/functional/environment tests/functional/playbook/ping.yml --user ubuntu
running: ansible_playbook --inventory tests/functional/environment/hosts tests/functional/playbook/ping.yml --user ubuntu

PLAY [ensure connectivity to all nodes] ****************************************

TASK [Gathering Facts] *********************************************************
ok: [default]

TASK [ansible setup] ***********************************************************
ok: [default]

PLAY RECAP *********************************************************************
default                    : ok=2    changed=0    unreachable=0    failed=0   
```