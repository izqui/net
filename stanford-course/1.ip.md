#IP
###Packet routing
IP is the network layer and it is the one in charge of delivering packets making best effort.

`$ traceroute -w 1 www.izqui.me`

###Packet switching
* Source routing: Specifying the whole route at the source. All the step until hitting the destination.
* Each router decides which is the next hop for getting to the destination.

A flow is a collection of datagrams belonging to the same end-to-end communication.

* They are simple. They forward packets independently.
* They are efficient. It is possible to share a link among many flows. 

###Layering
Design principle used in software architecture.

It allows you to abstract the details, and hides implementation details. This separation of concerns allows to build on top of the stack and be able to modify a part.

Each layer only communicate with the layers above and below. We can improve each layer independently.

* Modularity
* Well defined service
* Reuse
* Separation of concerns -> improvement

###IPv4
The goal was to stitch many different networks together.
32 bits long. Writen in 4 octets: [0-255].[0-255].[0-255].[0-255]

IP delivers packets to a device with that address.

The netmask specifies the part of the IP address that has to be the same for two computers to be on the same network, so a netmask like:

255.255.255.0 = {1111111}.{11111111}.{11111111}.0 Means that if the 3 first octets are the same, both computers are in the same network.

###Longest Prefix Match (IP Routing)
In each hop of each packet, the router needs to decide to what is going to be the next hop.

LPM is the algorithm routers use in order to forward packets:

default -> link 1
171.2.xxx.xxx -> link 3
171.2.34.xxx -> link 4

If there is a connection going to 171.2.34.11 it will go through link 4, because it is the longer match, the most specific one, even though it also matched the other ones.

###Address Resolution Protocol
Maps layer 2 addresses (link layer) with layer 3 addresses (network layer)

I'm on a network with netmask 255.255.255.0, me: 192.168.0.5 wants to send a packet to 192.168.0.1, but I don't know its link layer address (aka MAC address), so I broadcast a message to the network asking whether someone knows the address of that particular IP address. ARP requests are redundant, that means that when making a request you send your IP address and MAC address, so everyone that hears this can update his record. Reply to ARP request is only sent to the individual who made the request.

Assumption is that ARP records don't change very frequently. OSX updates it every 20 minutes, Cisco devices every 4 hours.



