package main

import (
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/bus"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/help"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/horario"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/menu"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/news"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/salas"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/status"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/subscription"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/tts"

	"github.com/guad/commander"
)

func route(cmd *commander.CommandGroup, cfg config) {
	cmd.Command("/menu", menu.CafeTodayCmd)
	cmd.Command("/menu2", menu.CafeTomorrowCmd)
	cmd.Command("/salas", salas.SalasCmd)
	cmd.Command("/noticias", news.NewsCmd)
	cmd.Command("/avisos", news.AvisosCmd)
	cmd.Command("/coredumped", news.CoreCmd)

	cmd.Command("/help", help.HelpCmd)
	cmd.Command("/start", help.HelpCmd)

	cmd.Command("/subscribe {feed?}", subscription.SubscribeCmd(cfg.db))
	cmd.Command("/unsubscribe {feed?}", subscription.UnsubscribeCmd(cfg.db))
	cmd.Command("/canales", subscription.GetAllChannelsCommand)

	go subscription.StartMonitoringSubscriptions(cfg.redis, cfg.bot, cfg.db)

	cmd.Command("/horario {grupo?}", horario.HorarioCmd(cfg.redis))
	cmd.Command("/horario2 {grupo?}", horario.HorarioWeekCmd(cfg.redis))

	cmd.Command("/status", status.StatusCmd)
	cmd.Command("/status_bot", status.BotStatusCmd)

	cmd.Command("/bus {stop:int?}", bus.BusCmd)

	cmd.Command("/tts", tts.TtsCmd)
}
