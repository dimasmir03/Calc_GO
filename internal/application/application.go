package application

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/dimasmir03/yandex-licey-go/1.13/pkg/rpn"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

// Функция запуска приложения
// тут будем читать введенную строку и после нажатия на ENTER писать результат работы программы на экран
// если пользователь ввел exit - то останавливаем приложение
func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		// убираем пробелы, чтобы оставить только вычисляемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("application was successfuly closed")
			return nil
		}
		// вычисляем выражение
		result, err := rpn.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed with error: ", err)
		} else {
			log.Println(text, "-", result)
		}
	}
}
