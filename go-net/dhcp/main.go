package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

// DHCPメッセージの定数（各DHCPメッセージタイプを定義）
const (
	DHCPDiscover = 1 // クライアントがIPアドレスを要求するときに送信される
	DHCPOffer    = 2 // サーバーがクライアントに提案するIPアドレスを送信
	DHCPRequest  = 3 // クライアントが提案されたIPアドレスを要求
	DHCPAck      = 5 // サーバーがIPアドレスの割り当てを確認
)

// DHCPパケットの基本構造（固定サイズ部分のフィールド定義）
type DHCPPacket struct {
	Op     byte      // メッセージの種類（1: Request, 2: Reply）
	Htype  byte      // ハードウェアタイプ（1: Ethernet）
	Hlen   byte      // ハードウェアアドレスの長さ（6: Ethernetアドレス長）
	Hops   byte      // 中継エージェントによるホップ数
	Xid    uint32    // トランザクションID（クライアントが生成するランダム値）
	Secs   uint16    // 経過時間（秒）
	Flags  uint16    // 各種フラグ
	Ciaddr [4]byte   // クライアントが既に持っているIPアドレス
	Yiaddr [4]byte   // サーバーが提供するクライアントのIPアドレス
	Siaddr [4]byte   // 次に使用するサーバーのIPアドレス
	Giaddr [4]byte   // 中継エージェントのIPアドレス
	Chaddr [16]byte  // クライアントのハードウェアアドレス
	Sname  [64]byte  // オプションのサーバーホスト名
	File   [128]byte // 起動ファイル名（オプション）
}

// UDPソケットを開き、クライアントからのリクエストを待つ
func main() {
	// UDPソケットをポート67でリッスン（DHCPサーバーが使用するポート）
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero, // すべてのインターフェースで待ち受ける
		Port: 67,           // DHCPサーバーが使用するポート番号
	})
	if err != nil {
		log.Fatalf("ソケット作成エラー: %v", err)
	}
	defer conn.Close()

	log.Println("DHCPサーバーがポート67で待ち受けています...")

	// 永遠にリクエストを待ち続けるループ
	for {
		handleDHCPRequest(conn)
	}
}

// DHCPリクエストの処理
func handleDHCPRequest(conn *net.UDPConn) {
	buffer := make([]byte, 1500) // 最大1500バイトのバッファを用意

	// クライアントからのメッセージを受信
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("データ受信エラー: %v", err)
		return
	}

	log.Printf("DHCPメッセージ受信: %dバイト from %v", n, addr)

	// 受信したデータをDHCPパケットとして解析
	packet, options, err := parseDHCPPacket(buffer[:n])
	if err != nil {
		log.Printf("パケット解析エラー: %v", err)
		return
	}

	// DHCPオプションからメッセージタイプを取得
	messageType := getDHCPMessageType(options)
	if messageType == 0 {
		log.Println("不明なDHCPメッセージタイプ")
		return
	}

	// メッセージタイプに応じた処理
	switch messageType {
	case DHCPDiscover:
		log.Println("DHCPDISCOVERを受信")
		sendDHCPOffer(conn, addr, packet) // DHCPDISCOVERに対してDHCPOFFERを送信
	case DHCPRequest:
		log.Println("DHCPREQUESTを受信")
		sendDHCPAck(conn, addr, packet)   // DHCPREQUESTに対してDHCPACKを送信
	}
}

// 受信したデータをDHCPパケットとして解析
func parseDHCPPacket(data []byte) (DHCPPacket, []byte, error) {
	var packet DHCPPacket
	reader := bytes.NewReader(data)

	// DHCPパケットの固定部分を読み込む（バイト列から構造体に変換）
	if err := binary.Read(reader, binary.BigEndian, &packet); err != nil {
		return packet, nil, fmt.Errorf("パケット読み込みエラー: %w", err)
	}

	// オプション部分を取り出す（DHCPパケットのヘッダー以降）
	options := data[240:] // 240バイト以降がオプションデータ

	return packet, options, nil
}

// DHCPオプションからメッセージタイプを取得
func getDHCPMessageType(options []byte) byte {
	for i := 0; i < len(options)-1; i++ {
		if options[i] == 53 && len(options) > i+1 {
			return options[i+1] // メッセージタイプの値を返す
		}
	}
	return 0 // 不明なタイプ
}

// DHCPOFFERメッセージをクライアントに送信
func sendDHCPOffer(conn *net.UDPConn, addr *net.UDPAddr, discover DHCPPacket) {
	var offer DHCPPacket

	// クライアントからの基本情報をコピー
	offer.Op = DHCPOffer         // メッセージタイプをDHCPOFFERに設定
	offer.Htype = discover.Htype // ハードウェアタイプをコピー
	offer.Hlen = discover.Hlen   // ハードウェアアドレスの長さをコピー
	offer.Xid = discover.Xid     // トランザクションIDをコピー
	copy(offer.Chaddr[:], discover.Chaddr[:]) // ハードウェアアドレスをコピー

	// 提供するIPアドレスを設定
	offer.Yiaddr = [4]byte{192, 168, 0, 100}

	// オプションの設定
	options := []byte{
		53, 1, DHCPOffer,            // DHCPメッセージタイプ
		1, 4, 255, 255, 255, 0,      // サブネットマスク
		54, 4, 192, 168, 0, 1,       // サーバー識別子
		51, 4, 0, 1, 81, 128,        // リース時間
		255,                         // オプション終了
	}

	// バイト列にパケットをエンコード
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, offer); err != nil {
		log.Printf("パケット作成エラー: %v", err)
		return
	}

	// オプションをバッファに追加
	buffer.Write(options)

	// クライアントにDHCPOFFERを送信
	if _, err := conn.WriteToUDP(buffer.Bytes(), addr); err != nil {
		log.Printf("DHCPOFFER送信エラー: %v", err)
	}
	log.Println("DHCPOFFERを送信しました")
}

// DHCPACKメッセージをクライアントに送信
func sendDHCPAck(conn *net.UDPConn, addr *net.UDPAddr, request DHCPPacket) {
	var ack DHCPPacket

	// クライアントからの基本情報をコピー
	ack.Op = DHCPAck            // メッセージタイプをDHCPACKに設定
	ack.Htype = request.Htype    // ハードウェアタイプをコピー
	ack.Hlen = request.Hlen      // ハードウェアアドレスの長さをコピー
	ack.Xid = request.Xid        // トランザクションIDをコピー
	copy(ack.Chaddr[:], request.Chaddr[:]) // ハードウェアアドレスをコピー

	// リクエストされたIPアドレスを設定
	ack.Yiaddr = request.Yiaddr

	// オプションの設定
	options := []byte{
		53, 1, DHCPAck,               // DHCPメッセージタイプ (53: メッセージタイプ, 1バイト, 値はDHCPAck)
		1, 4, 255, 255, 255, 0,       // サブネットマスク (1: サブネットマスクオプション, 4バイト, 255.255.255.0)
		54, 4, 192, 168, 0, 1,        // サーバー識別子 (54: サーバー識別子オプション, 4バイト, 192.168.0.1)
		51, 4, 0, 1, 81, 128,         // リース時間 (51: リース時間オプション, 4バイト, 3600秒=1時間)
		255,                          // オプション終了 (255: オプション終了)
	}

	// バイト列にパケットをエンコード
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, ack); err != nil {
		log.Printf("パケット作成エラー: %v", err)
		return
	}

	// オプションをバッファに追加
	buffer.Write(options)

	// クライアントにDHCPACKを送信
	if _, err := conn.WriteToUDP(buffer.Bytes(), addr); err != nil {
		log.Printf("DHCPACK送信エラー: %v", err)
	}
	log.Println("DHCPACKを送信しました")
}
