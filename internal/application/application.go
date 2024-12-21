package application

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dimasmir03/Calc_GO/pkg/calculation"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("port")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	// config для хранения порта приложения
	config *Config
	// sirupsen/logrus для логирования
	log *logrus.Logger
	// gorilla/mux для роутинга
	r *mux.Router
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
		log:    logrus.New(),
		r:      mux.NewRouter(),
	}
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
		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed with error: ", err)
		} else {
			log.Println(text, "-", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type SuccessResponse struct {
	Result string `json:"result"`
}

type FailedResponse struct {
	Error string `json:"error"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		if errors.Is(err, calculation.ErrInvalidExpression) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(&FailedResponse{Error: err.Error()})
		} else if errors.Is(err, calculation.ErrDivisionByZero) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(&FailedResponse{Error: err.Error()})
		} else if errors.Is(err, calculation.ErrInvalidCharacter) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(&FailedResponse{Error: err.Error()})
		} else if errors.Is(err, calculation.ErrInvalidToken) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(&FailedResponse{Error: err.Error()})
		} else if errors.Is(err, calculation.ErrMismatchParentheses) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(&FailedResponse{Error: err.Error()})
		} else if errors.Is(err, calculation.ErrUnknowOperator) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(&FailedResponse{Error: err.Error()})
		} else {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&FailedResponse{Error: http.StatusText(500)})
		}

	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&SuccessResponse{Result: fmt.Sprintf("%f", result)})
	}
}

func Logging(logger *logrus.Logger) mux.MiddlewareFunc {
	//middleware для логирования
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, req)
			logger.Infof("HTTP Запрос %s %s %s", req.Method, req.RequestURI, time.Since(start))
		})
	}
}

func CheckMethod(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		//middleware для проверки метода POST
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// проверка что запрос отправлен методов POST
			if req.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				json.NewEncoder(w).Encode(&FailedResponse{Error: ErrInvalidMethod.Error()})
				logger.Error(ErrInvalidMethod.Error())
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func (a *Application) RunServer() {
	//Использование middleware для логирования
	a.r.Use(Logging(a.log))
	a.r.Use(CheckMethod(a.log))
	// единственный endpoint приложения который принимает запрос только метода POST
	a.r.HandleFunc("/api/v1/calculate", CalcHandler)
	a.log.Infof("Starting server on :%s", a.config.Addr)
	if err := http.ListenAndServe(":"+a.config.Addr, a.r); err != nil {
		a.log.Errorf("Failed start server: %s", err.Error())
	}
}
