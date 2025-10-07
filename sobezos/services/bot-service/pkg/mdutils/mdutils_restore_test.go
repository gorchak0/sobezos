package mdutils

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var text = `/taskedit
$id
6
$tags
мод1
ОС
$question
Что такое сокет?
 $answer
Сокет — это программный интерфейс (endpoint) для обмена данными между процессами, который может работать как на одной машине, так и по сети.

Он связывает процесс с конечной точкой связи
🔹IP-адрес + порт для сетевых сокетов
🔹путь к файлу для локальных Unix сокетов.

Поддерживает двунаправленный обмен данными.

Sockets в Go:
func main() {

 listener, _ := net.Listen("tcp", ":8080") //💡

 defer listener.Close()
 fmt.Println("Listening on :8080")

 for {
  conn, _ := listener.Accept()
  go handleConnection(conn)
 }

}

func handleConnection(conn net.Conn) {
 defer conn.Close()
 fmt.Fprintf(conn, "Hello from Go server\n")
}`

var expected1 = `
$id
6
$tags
мод1
ОС
$question
Что такое сокет?
$answer
*Сокет* — это программный интерфейс (endpoint) для __обмена данными между процессами__, который может работать как на одной машине, так и по сети.

Он связывает процесс с конечной точкой связи
🔹IP-адрес + порт для сетевых сокетов
🔹путь к файлу для локальных Unix сокетов.

*Поддерживает двунаправленный обмен данными.*

__Sockets в Go:__
`

var expected2 = `
func main() {

 listener, _ := net.Listen("tcp", ":8080") //💡

 defer listener.Close()
 fmt.Println("Listening on :8080")

 for {
  conn, _ := listener.Accept()
  go handleConnection(conn)
 }

}

func handleConnection(conn net.Conn) {
 defer conn.Close()
 fmt.Fprintf(conn, "Hello from Go server\n")
}
`

var expected3 = expected1 + "```" + expected2 + "```"

func TestRestore(t *testing.T) {
	proc := NewMarkdownV2Processor()

	// Пример текста и entities

	entities := []tgbotapi.MessageEntity{
		{Type: "bot_command", Offset: 0, Length: 9},
		{Type: "bold", Offset: 66, Length: 5},
		{Type: "underline", Offset: 115, Length: 31},
		{Type: "bold", Offset: 335, Length: 43},
		{Type: "underline", Offset: 380, Length: 13},
		{Type: "pre", Offset: 394, Length: 302, Language: "go"},
	}

	entities2 := []tgbotapi.MessageEntity{
		{Type: "bot_command", Offset: 0, Length: 10},
	}

	result := proc.AddMD("/taskgetid 1", entities2)
	fmt.Printf("\n✅✅✅\nRestore result1: %s\n✅✅✅\n", result)

	result = proc.AddMD(text, entities)
	fmt.Printf("\n✅✅✅\nRestore result2: %s\n✅✅✅\n", result)
	fmt.Printf("\n✅✅✅\nExpected result: %s\n✅✅✅\n", expected3)
}
