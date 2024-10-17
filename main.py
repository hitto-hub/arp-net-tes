from scapy.all import ARP, Ether, srp1, conf

####################
# IPv4 ARP Request #
####################

# 送信先のIPv4アドレス
# DST_IP_V4 = "192.168.0.200"
DST_IP_V4 = "192.168.0.44"

# 使用するネットワークインターフェースを指定
conf.iface = "en0"  # 適切なインターフェース名を指定してください

def arp():
    # ARPリクエストパケットの作成
    pkt = Ether(dst="ff:ff:ff:ff:ff:ff") / ARP(op=1, pdst=DST_IP_V4)
    
    # ARPリクエストの送信とレスポンスの取得
    response = srp1(pkt, timeout=2, verbose=False)  # タイムアウトを2秒に設定
    
    # レスポンスがある場合
    if response:
        print(f"応答あり: {response[ARP].psrc}のMACアドレスは {response[Ether].src}")
    else:
        print("応答なし")

# ARPリクエストを実行
arp()
