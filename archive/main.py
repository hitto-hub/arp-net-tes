from scapy.all import IP, ICMP, TCP, Ether, UDP

a=IP(ttl=10)
print("a:", a)
# < IP ttl=10 |>
print("a.src:",a.src)
# ’127.0.0.1’
a.dst="192.168.1.1"
print(a)
# < IP ttl=10 dst=192.168.1.1 |>
print(a.src)
# ’192.168.8.14’
del(a.ttl)
print("a",a)
# < IP dst=192.168.1.1 |>
print("a.ttl",a.ttl)
# 64
print()
print(IP())
# <IP |>
print(IP()/TCP())
# <IP frag=0 proto=TCP |<TCP |>>
print(Ether()/IP()/TCP())
print(IP()/TCP()/"GET / HTTP/1.0\r\n\r\n")
# <IP frag=0 proto=TCP |<TCP |<Raw load='GET / HTTP/1.0\r\n\r\n' |>>>
print(Ether()/IP()/IP()/UDP())
# <Ether type=0x800 |<IP frag=0 proto=IP |<IP frag=0 proto=UDP |<UDP |>>>>
print(IP(proto=55)/TCP())
# <IP frag=0 proto=55 |<TCP |>>
