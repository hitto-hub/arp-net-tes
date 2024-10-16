import logging
from scapy.all import ARP, Ether, srp

# ログ設定
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', handlers=[
    logging.FileHandler("arp_log.txt", encoding='utf-8'),
    logging.StreamHandler()
])

def send_arp_request(target_ip, iface):
    try:
        logging.info(f"Sending ARP request to {target_ip} on interface {iface}")
        
        # ARPリクエストパケットを作成
        arp_request = ARP(pdst=target_ip)
        ether_frame = Ether(dst="ff:ff:ff:ff:ff:ff") / arp_request

        # パケットの詳細をログに記録
        logging.info(f"ARP request packet: {ether_frame.summary()}")
        
        # 指定したインターフェースにパケットを送信し、応答を待つ
        answered, unanswered = srp(ether_frame, iface=iface, timeout=2, verbose=False)
        
        if answered:
            for sent, received in answered:
                logging.info(f"Received ARP response from {received.psrc} (MAC: {received.hwsrc})")
        else:
            logging.info("No ARP response received")
    except PermissionError:
        logging.error("Error sending ARP request: Operation not permitted. Try running the script with elevated privileges.")
    except Exception as e:
        logging.error(f"Error sending ARP request: {e}")

if __name__ == "__main__":
    target_ip = "192.168.1.1"  # 送信先のIPアドレスを指定
    iface = "en0"  # 使用するNICのインターフェース名を指定 (適切なインターフェース名に変更)
    send_arp_request(target_ip, iface)
