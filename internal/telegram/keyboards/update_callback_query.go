package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetCRUDKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Crear", "crud_crear"),
			tgbotapi.NewInlineKeyboardButtonData("Llistar", "crud_llistar"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Actualitzar", "crud_actualitzar"),
			tgbotapi.NewInlineKeyboardButtonData("Eliminar", "crud_eliminar"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Tornar", "back"),
			tgbotapi.NewInlineKeyboardButtonData("Cancelar", "cancel"),
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
			tgbotapi.NewInlineKeyboardButtonData("Storage", fmt.Sprintf("get_storage_info_%s", instance)),
			tgbotapi.NewInlineKeyboardButtonData("Ports up", fmt.Sprintf("get_active_ports_%s", instance)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Tornar", "back_instance"),
		),
	)
}

func GetDiscoverKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Node Exporter", "discover_node_exporter"),
			tgbotapi.NewInlineKeyboardButtonData("Port Exporter", "discover_port_exporter"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Tornar", "back"),
		),
	)
}
