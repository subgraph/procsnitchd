
# procsnitchd

UNIX domain socket service providing process information for local network connections


# install

procsnitchd systemd unit file

```
[Unit]
Description=procsnitchd
After=tor.service

[Service]
Type=simple
User=root
Group=root
ExecStart=/usr/local/sbin/procsnitchd -socket /var/run/procsnitchd/rpc.socket -group procsnitch
```


Create the ``procsnitch`` group and directory.

```
groupadd --system procsnitch
mkdir /var/run/procsnitchd
chgrp procsnitch /var/run/procsnitchd
```


Programs using procsnitch should be in the ``procsnitch`` group so that they
can read/write the UNIX domain socket ``/var/run/procsnitchd/rpc.socket``.

