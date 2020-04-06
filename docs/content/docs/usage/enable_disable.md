---
title: Enable & Disable
weight: 20
---


## Enable a profile

You can enable any profile, the routing will react to it state. 

Command:

`hostctl enable -p awesome` 

Output:
```bash
+---------+--------+----------------+------------------------------+
| PROFILE | STATUS |       IP       |            DOMAIN            |
+---------+--------+----------------+------------------------------+
...
+---------+--------+----------------+------------------------------+
| awesome | on     | 127.0.0.1      | web.my-awesome-project.local |
| awesome | on     | 127.0.0.1      | api.my-awesome-project.local |
+---------+--------+----------------+------------------------------+
```


## Disable a profile

You can disable any profile, all routing for that profile will stop working. 

{{<info>}}
Disabling a profile does not remove the content from the hosts file, this way you can re-enable it later on.
{{</info>}}

Command:

`hostctl disable -p awesome` 

Output:
```bash
+---------+--------+----------------+------------------------------+
| PROFILE | STATUS |       IP       |            DOMAIN            |
+---------+--------+----------------+------------------------------+
...
+---------+--------+----------------+------------------------------+
| awesome | off    | 127.0.0.1      | web.my-awesome-project.local |
| awesome | off    | 127.0.0.1      | api.my-awesome-project.local |
+---------+--------+----------------+------------------------------+
```
