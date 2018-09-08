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
        if update.message.chat_id < 0:  # This pre-check is necessary if we do not want to spam the logs with "BadRequest: Message can't be deleted" as this bot has no power to remove user messages in private chats.
            logger.exception("handle malformed requests - read more below!")
    except TimedOut:
        logger.exception("handle slow connection problems")
    except NetworkError:
        logger.exception("handle other connection problems")
    except ChatMigrated as e:
        logger.exception("the chat_id of a group has changed, use " + e.new_chat_id + " instead")
    except TelegramError:
        logger.exception("There is some error with Telegram")

def is_admin(user_id):
    if user_id in settings.admin_ids:
        return True
    return False

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


def delete_message(bot, update):
    bot.deleteMessage(update.message.chat_id, update.message.message_id)


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
            last_function_calls[name][chat_id] = now
            return True
    else:
        last_function_calls[name] = {chat_id: now}
        return True


def help_command(bot, update):
    if is_call_available("help_command", update.message.chat_id, 180):
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
    if is_admin(update.message.from_user.id):
        logger.info("Reloading settings")
        load_settings()
        bot.send_message(chat_id=update.message.chat_id, text="Datos cargados")
    delete_message(bot, update)


def news_command(bot, update):
    if is_call_available("news_command", update.message.chat_id, 180):
        log_message(update)
        logger.info("Getting news")
        text = "Estas son las últimas noticias que aparecen en la web: \n"
        news_list = news_json_scraper()
        for idx, new in enumerate(news_list):
            text = text + str(idx) + ") " + news_list[new]["a-link"] + "\n"
        bot.sendMessage(chat_id=update.message.chat.id, text=text, parse_mode=telegram.ParseMode.HTML)


def events_command(bot, update):
    if is_call_available("events", update.message.chat_id, 180):
        log_message(update)
        logger.info("Getting news")
        text = "Estas son los últimos eventos que aparecen en la web: \n"
        events_list = events_json_scraper()
        for idx, new in enumerate(events_list):
            text = text + str(idx) + ") " + events_list[new]["a-link"] + "\n"
        bot.sendMessage(chat_id=update.message.chat.id, text=text, parse_mode=telegram.ParseMode.HTML)


def notifications_command(bot, update):
    if is_call_available("help", update.message.chat_id, 180):
        log_message(update)
        logger.info("Getting news")
        text = "Estas son los últimos avisos que aparecen en la web: \n"
        notifications_list = avisos_json_scraper()
        for idx, new in enumerate(notifications_list):
            text = text + str(idx) + ") " + notifications_list[new]["a-link"] + "\n"
        bot.sendMessage(chat_id=update.message.chat.id, text=text, parse_mode=telegram.ParseMode.HTML)
    else:
        delete_message(bot, update)


def schedule_command(bot, update, args):  # Add arguments for checking other's group schedule
    global chat_ids_list
    global schedule_list

    def schedule_parser(schedule):
        parsed_schedule = [""]
        schedule_keys = sorted(schedule, key=lambda s: int(s.split(":")[0]))
        for hour in schedule_keys:
            parsed_schedule.append("A las %sh -> %s" % (hour, schedule[hour]))
        return "\n".join(parsed_schedule)

    if is_call_available("schedule", update.message.chat_id, 20):
        log_message(update)
        try:
            if update.message.chat_id < 0:  # ID's below 0 are groups.
                group = update.message.chat.title.replace(" ETSISI", "")  # Borro contenido de los títulos que me sobra
                if args:  # Ignore arguments if call is recieved from group.
                    bot.send_message(chat_id=update.message.chat_id,
                                     text="No respondo peticiones de horario de otros grupos aquí para evitar SPAM. Inicia un chat privado conmigo y pregúntame.")
                    return
            else:
                group = args[0].upper()
            if (str(datetime.datetime.today().weekday()) in range(1, 5)):
                text = schedule_parser(schedule_list[group][str(datetime.datetime.today().weekday())])
                text = "Horario de hoy para " + group + ":" + text
                bot.send_message(chat_id=update.message.chat_id, text=text)
            else:
                bot.send_message(chat_id=update.message.chat_id, text="Hoy no hay horario")



        except:
            bot.send_message(chat_id=update.message.chat_id, text="No he podido procesar tu solicitud de horario.")
    else:
        delete_message(bot, update)


if __name__ == "__main__":
    print("ETSISI Bot: Starting...")

    logger = get_logger("bot_starter", True)
    load_settings()

    try:
        logger.info("Conectando con la API de Telegram.")
        updater = Updater(settings.telegram_token)
        dispatcher = updater.dispatcher
        dispatcher.add_handler(CommandHandler('help', help_command))
        dispatcher.add_handler(CommandHandler('reload', reload_data))
        dispatcher.add_handler(CommandHandler('noticias', news_command))
        dispatcher.add_handler(CommandHandler('eventos', events_command))
        dispatcher.add_handler(CommandHandler('avisos', notifications_command))
        dispatcher.add_handler(CommandHandler('horario', schedule_command, pass_args=True))
        dispatcher.add_handler(MessageHandler(Filters.status_update, delete_message))
        dispatcher.add_error_handler(error_callback)

    except Exception as ex:
        logger.exception("Error al conectar con la API de Telegram.")
        quit()

    updater.start_polling()
    logger.info("ETSISI Bot: Estoy escuchando.")
    updater.idle()
