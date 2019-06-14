package main

import (
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/bus"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/exam"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/guides"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/help"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/horario"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/janitor"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/menu"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/news"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/salas"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/status"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/stub"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/subscription"
	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/cmds/tts"

	"github.com/guad/commander"
)

// Note that underscores (_) are forbidden for command names.
func route(cmd *commander.CommandGroup, cfg config, callbacks *commander.CommandGroup) {
	cmd.Command("/menu", menu.CafeTodayCmd)
	cmd.Command("/menu2", menu.CafeTomorrowCmd)
	cmd.Command("/salas", salas.SalasCmd)
	cmd.Command("/noticias", news.NewsCmd)
	cmd.Command("/avisos", news.AvisosCmd)
	cmd.Command("/coredumped", news.CoreCmd)

	cmd.Command("/help", help.HelpCmd)
	cmd.Command("/start {data*}", help.Start)

	cmd.Command("/subscribe {feed?}", subscription.SubscribeCmd(cfg.db))
	cmd.Command("/unsubscribe {feed?}", subscription.UnsubscribeCmd(cfg.db))
	cmd.Command("/canales", subscription.GetAllChannelsCommand)

	go subscription.StartMonitoringSubscriptions(cfg.redis, cfg.bot, cfg.db)

	cmd.Command("/horario {grupo?}", horario.HorarioCmd)
	cmd.Command("/horario2 {grupo?}", horario.HorarioWeekCmd)

	cmd.Command("/status", stub.Middleware(status.StatusCmd))
	cmd.Command("/statusbot", stub.Middleware(status.BotStatusCmd))

	cmd.Command("/bus {stop:int?}", bus.BusCmd)

	cmd.Command("/tts", tts.TtsCmd)

	cmd.Command("/exam", exam.ExamCmd)

	cmd.Command("/guias", guides.GuideCmd)
	cmd.Command("/gg {code}", guides.DownloadGuideCmd)

	cmd.Command("/manage", janitor.AdminOnlyMiddleware(janitor.Manage))
	cmd.Command("/ban",
		janitor.AdminOnlyMiddleware(
			janitor.ManagedOnlyMiddleware(janitor.Ban),
		),
	)

	cmd.Command("/unban {user:int}",
		janitor.AdminOnlyMiddleware(
			janitor.ManagedOnlyMiddleware(janitor.Unban),
		),
	)

	// Callbacks
	callbacks.Command("/gpag {grado} {offset:int}", guides.PaginateGradoCallback)

	callbacks.Command("/exyear {grado}", exam.SelectYearCb)
	callbacks.Command("/exshow {grado} {curso:int}", exam.ShowExamsCb)

	callbacks.Command("/jannyrefresh {chatid} {public}",
		janitor.AdminOnlyMiddleware(
			janitor.ManagedOnlyMiddleware(janitor.RefreshCb),
		),
	)

	callbacks.Command("/janmnypublictoggle {chatid} {public}",
		janitor.AdminOnlyMiddleware(
			janitor.ManagedOnlyMiddleware(janitor.TogglePublicCb),
		),
	)

	callbacks.Command("/jannydisable {chatid}",
		janitor.AdminOnlyMiddleware(
			janitor.ManagedOnlyMiddleware(janitor.DisableCb),
		),
	)

	// Events
	cmd.Event("text", janitor.OnMessage)
}
