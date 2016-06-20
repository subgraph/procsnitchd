
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
groupadd procsnitch
mkdir /var/run/procsnitchd
chown procsnitch:procsnitch /var/run/procsnitchd
```


Be in the ``procsnitch`` group so that you can read/write the UNIX domain sock ``/var/run/procsnitchd/rpc.socket``.

