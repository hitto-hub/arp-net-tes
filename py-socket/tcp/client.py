import socket

def start_client():
    # TCP/IPソケットを作成
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # サーバーのアドレスとポートに接続
    client_socket.connect(('localhost', 8080))

    # サーバーにメッセージを送信
    message = "こんにちは、サーバー！"
    client_socket.sendall(message.encode('utf-8'))

    # サーバーからの応答を受信
    response = client_socket.recv(1024)
    print(f"サーバーからの応答: {response.decode('utf-8')}")

    # 接続を閉じる
    client_socket.close()

# クライアントを開始
start_client()
