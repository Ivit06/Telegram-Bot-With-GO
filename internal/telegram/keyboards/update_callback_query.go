package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetCRUDKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("CREAR", "crud_crear"),
			tgbotapi.NewInlineKeyboardButtonData("LLISTAR", "crud_llistar"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ACTUALITZAR", "crud_actualitzar"),
			tgbotapi.NewInlineKeyboardButtonData("ELIMINAR", "crud_eliminar"),
		),
	)
}

func GetNodeMetricsKeyboard(instance string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("CPU", fmt.Sprintf("get_cpu_info_%s", instance)),
			tgbotapi.NewInlineKeyboardButtonData("RAM", fmt.Sprintf("get_ram_info_%s", instance)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("STORAGE", fmt.Sprintf("get_storage_info_%s", instance)),
			tgbotapi.NewInlineKeyboardButtonData("NETWORK", fmt.Sprintf("get_network_info_%s", instance)),
		),
	)
}

func GetDiscoverKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Node Exporter", "discover_node_exporter"),
			tgbotapi.NewInlineKeyboardButtonData("Port Exporter", "discover_port_exporter"),
		),
	)
}
