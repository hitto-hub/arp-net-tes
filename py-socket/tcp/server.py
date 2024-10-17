import socket

def start_server():
    # TCP/IPソケットを作成
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # サーバーをローカルホストのポート8080でバインド
    server_socket.bind(('localhost', 8080))

    # クライアントからの接続を待ち受ける（最大5接続を待機可能）
    server_socket.listen(5)
    print("サーバーがポート8080で待ち受けています...")

    while True:
        # クライアントからの接続を受け入れる
        client_socket, client_address = server_socket.accept()
        print(f"クライアント {client_address} が接続しました")

        # クライアントからのデータを受信
        data = client_socket.recv(1024)
        if data:
            print(f"受信: {data.decode('utf-8')}")
            # クライアントに応答を送信
            client_socket.sendall("サーバーからの応答です".encode('utf-8'))

        # 接続を閉じる
        client_socket.close()

# サーバーを開始
start_server()
