# gosible

Gosbile is a wrapper around Ansible designed to implement a stronger "Infrastructure as Code" mentality.

When building and developing [Bluebox Cloud](https://github.com/blueboxgroup/ursula) and [Cuttle](https://github.com/ibm/cuttle) we found that
we had hundreds of deployments that were very similar so we looked to optimize our deployment workflow to solve that. We built a tool called [ursula-cli](https://github.com/blueboxgroup/ursula-cli) in Python, this is a rewrite of it in Golang. It supports most of our commonly used features, and some more.

When running a playbook with Ansible it keys off your inventory file and then uses its path to try and find `group_vars`, `host_vars`, etc. For Infrastructure as Code to work well you need to keep both your Ansible playbooks and your inventory in source control.  However we wanted to take a step further and be able to keep an ssh config, known hosts file, etc and have Ansible be aware of those things.

To do this we wrote a wrapper around Ansible that keys off an environment directory rather than the inventory file, and we keep all environments in git including their `ssh_config`, `ssh_known_hosts` and often a vagrantfile, heat template or ansible playbook for infrastructure provisioning.  Gosible looks for these and if they exist tells Ansible to use them. This also means as an Operator I always know how to SSH into a given environment `$ ssh -F infra/ssh_config elk01`.  It will also if prompted utilize a `requirements.txt` and a `virtualenv` directory in each environment to ensure the right version of ansible (and any other deps like boto or softlayer libs) are installed.

If you have a virtualenv directory in your current working directory or in the specified environment gosible will add `virtualenv/bin` to your PATH and attempt to use it to run ansible commands.

This pattern is what allows a surprisingly small team deploy and operate hundreds of [Openstack](https://github.com/blueboxgroup/ursula) deployments, and their [SRE Operations Platform](https://github.com/IBM/cuttle).

```
.
├── ansible.cfg
├── defaults.yml
├── bastion
│   ├── group_vars
│   │   ├── all.yml
│   ├── hosts
│   ├── host_vars
│   ├── ssh_config
│   ├── ssh_known_hosts
│   └── vagrant.yml
├── infra
│   ├── group_vars
│   │   ├── all.yml
│   ├── heat_stack.yml
│   ├── hosts
│   ├── host_vars
│   ├── requirements.txt
│   ├── ssh_config
│   ├── ssh_known_hosts
│   ├── vagrant.yml
│   └── virtualenv
```

## Install

While we're early on in development the easiest way to install Gosible is by running the following:

```
$ go get github.com/paulczar/gosible
```

You can also check for [prebuilt binary releases](https://github.com/paulczar/gosible/releases).

## Example Usage

Gosbile has some test environments suitable for demonstrating its usage.

```
$ git clone git@github.com:paulczar/gosible.git
$ make
$ bin/gosible
Gosible is a CLI tool designed to implement stronger 
Infrastructure-as-Code abilities to Ansible.

https://github.com/paulczar/gosible

Usage:
  gosible [command]
```

### Playbooks

To run a playbook across an environment we run `gosible playbook run --environment=/path/to/env playbook.yml`.
See `gosible playbook run --help` for a comprehensive set of arguments that can be set.  Any arguments set after the
playbook file will be passed onto ansible as is.

The following uses Gosible's `vagrant` provisioner that can be used to easily create a development environments on your local machine.

```bash
$ bin/gosible playbook run --provisioner=vagrant -e tests/functional/environment tests/functional/playbook/ping.yml --become

PLAY [ensure connectivity to all nodes] ****************************************

TASK [Gathering Facts] *********************************************************
ok: [default]

TASK [ansible setup] ***********************************************************
ok: [default]

PLAY RECAP *********************************************************************
default                    : ok=2    changed=0    unreachable=0    failed=0   
```

### Install Ansible

Gosible understands how to install Ansible:

```
$ gosible help install

Gosible install will install ansible and/or any dependencies provided.
If you do not specify a virtualenv or requirements file it will attempt
to find them in your local path, or it will try to find them in your
environment if you provide one.

Usage:
  gosible install [flags]

Flags:
  -e, --environment string    path to your gosible environment
  -h, --help                  help for install
  -r, --requirements string   path to requirements.txt
  -v, --virtualenv string     Path to VirtualEnv to use
```

You can see it in action by using the example environment in `tests/funtional/environment`

```
$ gosible install -e tests/functional/environment
checking if tests/functional/environment/requirements.txt exists... yes
Using VirtualEnv: tests/functional/environment/virtualenv
Using Pip: tests/functional/environment/virtualenv/bin/pip
==> Running: tests/functional/environment/virtualenv/bin/pip install -r tests/functional/environment/requirements.txt
Requirement already satisfied: ansible in ./tests/functional/environment/virtualenv/lib/python2.7/site-packages (from -r tests/functional/environment/requirements.txt (line 1))
Requirement already satisfied: softlayer in ./tests/functional/environment/virtualenv/lib/python2.7/site-packages (from -r tests/functional/environment/requirements.txt (line 2))
Requirement already satisfied: PyYAML in ./tests/functional/environment/virtualenv/lib/python2.7/site-packages (from ansible->-r tests/functional/environment/requirements.txt (line 1))
...
...
```

### Ping

Sometimes you just want to check if all the hosts in your environment are contactable, If you do not specify an environment
Gosible will attempt to use your current working directory as the environment. So we can do the following:

```
$ cd tests/functional/environment
$ ../../../bin/gosible ping
default | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
```

### Adhoc Commands

Sometimes you just want to run an adhoc command over all your hosts:

```
$ ../../../bin/gosible adhoc "hostname"                              
default | SUCCESS | rc=0 >>
ubuntu-xenial
```