#!/usr/bin/env python
# -*- coding: utf-8 -*-
import io
import json

import telegram
import datetime
import time
import os

from telegram.ext.jobqueue import Days

import logger
from logger import get_logger
from data_loader import DataLoader
import sys
from telegram.ext import Updater, CommandHandler, MessageHandler, BaseFilter, Filters
from random import normalvariate
from telegram.error import (TelegramError, Unauthorized, BadRequest,
                            TimedOut, ChatMigrated, NetworkError)
from etsisi_web_scraper import news_json_scraper, events_json_scraper, avisos_json_scraper
from upm_json_consultor import get_etsisi_degrees_info_json

reload(sys)
sys.setdefaultencoding('utf8')


def error_callback(bot, update, error):
    try:
        raise error
    except Unauthorized:
        logger.exception("remove update.message.chat_id from conversation list")
    except BadRequest:
        logger.exception("handle malformed requests - read more below!")
    except TimedOut:
        logger.exception("handle slow connection problems")
    except NetworkError:
        logger.exception("handle other connection problems")
    except ChatMigrated as e:
        logger.exception("the chat_id of a group has changed, use " + e.new_chat_id + " instead")
    except TelegramError:
        logger.exception("There is some error with Telegram")


def get_schedule():
    with io.open('horarios.json', 'r', encoding='utf8') as data_file:
        return json.load(data_file)


def get_chat_ids():
    with io.open('chat_ids.json', 'r', encoding='utf8') as data_file:
        return json.load(data_file)


def load_settings():
    global settings
    global last_function_calls
    global schedule_list
    global chat_ids_list
    settings = DataLoader()
    schedule_list = get_schedule()
    chat_ids_list = get_chat_ids()
    last_function_calls = {}


def is_call_available(name, chat_id, cooldown):
    global last_function_calls
    now = datetime.datetime.now()
    cooldown_time = datetime.datetime.now() - datetime.timedelta(minutes=cooldown)
    if name in last_function_calls.keys():
        if chat_id in last_function_calls[name].keys():
            if last_function_calls[name][chat_id] > cooldown_time:
                last_function_calls[name][chat_id] = now
                return False
            else:
                last_function_calls[name][chat_id] = now
                return True
        else:
            last_function_calls[name] = {chat_id: now}
            return True
    else:
        last_function_calls[name] = {chat_id: now}
        return True


def help(bot, update):
    log_message(update)
    bot.sendMessage(update.message.chat_id, settings.help_string, parse_mode=telegram.ParseMode.HTML)


def log_message(update):
    logger.info("He recibido: \"" + update.message.text + "\" de " + update.message.from_user.username + " [ID: " + str(
        update.message.chat_id) + "]")


def human_texting(string):
    wait_time = len(string) * normalvariate(0.1, 0.05)
    if wait_time > 8:
        wait_time = 8
    time.sleep(wait_time)


def reload_data(bot, update):
    if update.message.from_user.id == 15360527:
        logger.info("Reloading settings")
        load_settings()
        bot.send_message(chat_id=update.message.chat_id, text="Datos cargados")


def news(bot, update):
    logger.info("Getting news")
    text = "Estas son las últimas noticias que aparecen en la web: \n"
    news_list = news_json_scraper()
    for idx, new in enumerate(news_list):
        text = text + str(idx) + ") " + news_list[new]["a-link"] + "\n"
    bot.sendMessage(chat_id=update.message.chat.id, text=text, parse_mode=telegram.ParseMode.HTML)


def events(bot, update):
    logger.info("Getting news")
    text = "Estas son los últimos eventos que aparecen en la web: \n"
    events_list = events_json_scraper()
    for idx, new in enumerate(events_list):
        text = text + str(idx) + ") " + events_list[new]["a-link"] + "\n"
    bot.sendMessage(chat_id=update.message.chat.id, text=text, parse_mode=telegram.ParseMode.HTML)


def notifications(bot, update):
    logger.info("Getting news")
    text = "Estas son los últimos avisos que aparecen en la web: \n"
    notifications_list = avisos_json_scraper()
    for idx, new in enumerate(notifications_list):
        text = text + str(idx) + ") " + notifications_list[new]["a-link"] + "\n"
    bot.sendMessage(chat_id=update.message.chat.id, text=text, parse_mode=telegram.ParseMode.HTML)


def schedule(bot, update):
    global chat_ids_list
    global schedule_list

    def schedule_parser(schedule):
        print(schedule)
        parsedSchedule = [""]
        scheduleKeys = sorted(schedule, key=lambda s: int(s.split(":")[0]))
        for hour in scheduleKeys:
            parsedSchedule.append("A las %sh -> %s" % (hour, schedule[hour]))
        return "\n".join(parsedSchedule)

    try:
        group = update.message.chat.title.replace(" ETSISI", "")  # Borro contenido de los títulos que me sobra
        text = "Horario de hoy para " + group + ":" + schedule_parser(
            schedule_list[group][str(datetime.datetime.today().weekday())]) + "\n\n Gracias a Yadkee por su ayuda"
        bot.send_message(chat_id=update.message.chat_id, text=text)
    except:
        bot.send_message(chat_id=update.message.chat_id, text="No he podido procesar tu solicitud de horario.")


def delete_message(bot, update):
    bot.deleteMessage(update.message.chat_id, update.message.message_id)


if __name__ == "__main__":
    print("ETSISI Bot: Starting...")

    logger = get_logger("bot_starter", True)
    load_settings()

    try:
        logger.info("Conectando con la API de Telegram.")
        updater = Updater(settings.telegram_token)
        dispatcher = updater.dispatcher
        dispatcher.add_handler(CommandHandler('help', help))
        dispatcher.add_handler(CommandHandler('reload', reload_data))
        dispatcher.add_handler(CommandHandler('noticias', news))
        dispatcher.add_handler(CommandHandler('eventos', events))
        dispatcher.add_handler(CommandHandler('avisos', notifications))
        dispatcher.add_handler(CommandHandler('horario', schedule))
        dispatcher.add_handler(MessageHandler(Filters.status_update, delete_message))
        dispatcher.add_error_handler(error_callback)

    except Exception as ex:
        logger.exception("Error al conectar con la API de Telegram.")
        quit()

    updater.start_polling()
    logger.info("ETSISI Bot: Estoy escuchando.")
    updater.idle()
