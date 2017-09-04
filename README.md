# gosible

Gosbile is a wrapper around Ansible designed to implement a stronger "Infrastructure as Code" mentality.

Gosbile assumes you have a directory that describes the `environment` that you wish to run ansible over.  This
`environment` path contains at minimum your inventory (`./hosts`) along with any variables (`./group_vars/*`, `./host_vars/*`, etc).
It can also include an ssh config file `./ssh_config`, and an ssh known hosts file `./ssh_known_hosts`.

It is loosely based off the Bluebox project [ursula-cli](https://github.com/blueboxgroup/ursula-cli) which did much of the same thing and is used to deploy and manage hundreds of [Openstack](https://github.com/blueboxgroup/ursula) deployments, and their [operations platforms](https://github.com/IBM/cuttle).

The goal is to allow you to store everything involved in creating and deploying your infrastructure in git right alongside your Ansible Playbooks/Roles/etc.



## Install

While we're early on in development the easiest way to install Gosible is by running the following:

```
$ go get github.com/paulczar/gosible
```

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