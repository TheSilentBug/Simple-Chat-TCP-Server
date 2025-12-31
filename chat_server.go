package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var clients []net.Conn // لیست کلاینت‌های متصل

// تابعی برای ارسال پیام به تمام کلاینت‌ها
func broadcastMessage(message string, sender net.Conn) {
	for _, client := range clients {
		if client != sender { // پیام را به جز فرستنده به همه ارسال کن
			_, err := client.Write([]byte(message))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}
}

// تابعی برای مدیریت هر کلاینت
func handleClient(client net.Conn) {
	defer client.Close()

	// افزودن کلاینت به لیست
	clients = append(clients, client)

	scanner := bufio.NewScanner(client)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Received:", message) // چاپ پیام دریافتی

		if message == "exit" {
			break
		}

		// ارسال پیام به دیگر کلاینت‌ها
		broadcastMessage(message+"\n", client)
	}

	// حذف کلاینت از لیست در هنگام قطع ارتباط
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// تابعی برای راه‌اندازی سرور
func startServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Server started on port 8080...")

	// پذیرش اتصال‌ها
	for {
		client, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(client) // هر کلاینت در یک گوروتین مدیریت می‌شود
	}
}

func main() {
	startServer() // شروع سرور
}
