package mailing

import (
	"fmt"
	"strings"

	"github.com/nurtai325/kaspi-service/internal/models"
)

const (
	reviewLinkBase = "https://kaspi.kz/shop/review/productreview?rating=5"
)

func newOrderMessage(client models.Client, order models.Order) string {
	text := `Қайырлы күн, %s!
%s магазиніне тапсырыс бергеніңіз үшін рақмет
Сіздің тапсырысыңыз: %s
Жеткізу көрсетілген күні жүзеге асырылады.
Тапсырыс нөмірі: %s
Жұмыс уақыты: 10:00-20:00 дейын
Сұрақтар бойнша осы %s номерге хабарласыңыз

Тапсырысыңызды тез арада жинап, сізге жібереміз.
Саудаңыз сәтті болсын!


*

Добрый день, %s!
Спасибо за заказ в магазине %s
Вы заказали: %s.
Доставка будет осуществлена в указанную дату.
Номер заказа: %s
График работы с 10:00 до 20:00
По всем вопросам обращайтесь по %s этому номеру 
В ближайшее время мы соберём заказ и отправим вам.
Хороших покупок!`
	message := fmt.Sprintf(
		text,
		order.Customer,
		client.Name,
		parseEntries(order.Entries, true),
		order.Code,
		"+"+client.Phone,
		order.Customer,
		client.Name,
		parseEntries(order.Entries, false),
		order.Code,
		"+"+client.Phone,
	)
	return message
}

func completedOrderMessage(clientName string, order models.Order) string {
	reviewLink := fmt.Sprintf("%s&orderCode=%s&productCode=%s", reviewLinkBase, order.Code, order.ProductCode)
	text := `Қайырлы күн, %s!
%s дүкенінен сатып алуыңызбен құттықтаймыз!
Сізге бәрі ұнады деп үміттенеміз.
Сілтеме арқылы өтіп, біздің дүкеннің атауын көрсете отырып, пікір қалдыра аласыз ба, бұл біз үшін маңызды⤵️:
%s
Сілтеме белсенді болуы үшін жауап ретінде бірдеңе жазыңыз.

🔸🔸🔸

Добрый день, %s!
Поздравляем с покупкой с магазина %s!
Мы надеемся, что вам все понравилось.
Если вам не сложно, пожалуйста оставьте отзыв с указанием названия нашего магазина перейдя по ссылке⤵️:
%s
Чтобы ссылка стала активной, напишите пожалуйста что-нибудь в ответ.`
	message := fmt.Sprintf(text, order.Customer, clientName, reviewLink, order.Customer, clientName, reviewLink)
	return message
}

func parseEntries(entries []models.Entry, kz bool) string {
	parsed := strings.Builder{}
	length := len(entries)

	for i, entry := range entries {
		parsedEntry := ""
		if kz {
			parsedEntry = fmt.Sprintf(
				"%s, %d дана, бағасы: %d",
				entry.ProductName,
				entry.Quantity,
				int(entry.Price),
			)
		} else {
			parsedEntry = fmt.Sprintf(
				"%s, количество: %d, цена: %d",
				entry.ProductName,
				entry.Quantity,
				int(entry.Price),
			)
		}
		parsed.WriteString(parsedEntry)
		if i+1 != length {
			parsed.WriteString(",")
		}
	}
	parsed.WriteString(".")

	return parsed.String()
}
