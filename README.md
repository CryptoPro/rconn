rconn (r[everse] conn[ection]) is a multiplatform program for creating reverse connections. It lets you consume services that are behind NAT and/or firewall without adding firewall rules or port-forwarding. This is achieved by creating a connection from the node behind the firewall/NAT to a port on your local machine, and then a port is exposed in your machine through which you can connect to the service that is behind firewall/NAT. All traffic is routed through the initial connection that was opened by the machine behind firewall/NAT.

### Building
Build with: `go build`.

### Explanation
![diagram](https://github.com/jafarlihi/rconn/blob/master/diagram.png?raw=true)

Say your IP address is 11.11.11.11, and you've got machine 22.22.22.22 behind firewall/NAT and you want to connect to it via RDP. First you'd have to make sure your RDP server is running, normally on 3389. Now the problem is you can't connect to 3389 from outside because of NAT or firewall. Then in your local machine you'd have to run this:
`rconn -s 1111 2222`.
And in the machine behind firewall/NAT you'd have to run this:
`rconn -c 11.11.11.11 1111 127.0.0.1 3389`
Now you can connect to your own port 2222 with an RDP client, this will effectively be same as connecting to 22.22.22.22:3389. You can also change "127.0.0.1" to any other host the target machine has access to.

Usually most firewalls allow all outbound traffic but if this is not the case then you can try 80 or 443 instead of 1111.
