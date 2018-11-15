#!/usr/bin/env python
# -*- coding: utf-8 -*-
import io
import json

import telegram
import datetime
import time
import traceback
import re
import requests


import logger
from logger import get_logger
from data_loader import DataLoader
import sys
from telegram.ext import Updater, CommandHandler, MessageHandler, BaseFilter, Filters
from random import normalvariate
from telegram.error import (TelegramError, Unauthorized, BadRequest,
                            TimedOut, ChatMigrated, NetworkError)
from etsisi_web_scraper import news_json_scraper, events_json_scraper, avisos_json_scraper
from selenium import webdriver

reload(sys)
sys.setdefaultencoding('utf8')


def error_callback(bot, update, error):
    try:
        raise error
    except Unauthorized:
        logger.exception("remove update.message.chat_id from conversation list")
    except BadRequest as e:
        if update.message.chat_id < 0 and e == "Message can't be deleted":  # This pre-check is necessary if we do not want to spam the logs with "BadRequest: Message can't be deleted" as this bot has no power to remove user messages in private chats.
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


weekdays = ["lunes", "martes", "miércoles", "jueves", "viernes", "sábado", "domingo"]

codedays = {
    "L": 0,
    "M": 1,
    "X": 2,
    "J": 3,
    "V": 4,
    "S": 5,
    "D": 6,
    "LUNES": 0,
    "MARTES": 1,
    "MIÉRCOLES": 2,
    "MIERCOLES": 2,
    "JUEVES": 3,
    "VIERNES": 4,
    "SÁBADO": 5,
    "SABADO": 6,
    "DOMINGO": 7,
}


def get_schedule():
    with io.open('data/horarios.json', 'r', encoding='utf8') as data_file:
        return json.load(data_file)


def get_teachers():
    with io.open('data/profesores.json', 'r', encoding='utf8') as data_file:
        return json.load(data_file)


def get_chat_ids():
    with io.open('data/chat_ids.json', 'r', encoding='utf8') as data_file:
        return json.load(data_file)


def load_settings():
    global settings
    global last_function_calls
    global schedule_list
    global teacher_list
    global chat_ids_list
    settings = DataLoader()
    schedule_list = get_schedule()
    teacher_list = get_teachers()
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


def reset_call(name, chat_id):
    global last_function_calls
    reset_date = datetime.datetime.now() - datetime.timedelta(days=1)
    last_function_calls[name][chat_id] = reset_date


def help_command(bot, update):
    if is_call_available("help_command", update.message.chat_id, 180):
        log_message(update)
        bot.sendMessage(update.message.chat_id, settings.help_string, parse_mode=telegram.ParseMode.HTML)


def log_message(update):
    try:
        username = update.message.from_user.username
    except:
        username = "desconocido"
    try:
        text = update.message.text
    except:
        text = "algo"
    logger.info("He recibido: \"" + text + "\" de " + username + " [ID: " + str(
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


def calendar_command(bot, update):
    if is_call_available("calendario", update.message.chat_id, 180):
        log_message(update)
        bot.sendMessage(update.message.chat_id, settings.calendar_string, parse_mode=telegram.ParseMode.HTML)


def menu_command(bot, update):
    if is_call_available("cafe", update.message.chat_id, 180):
        text = ""
        log_message(update)

        try:
            data = requests.get("https://cafe.kolhos.chichasov.es/today").json()

            if "message" in data:  # Error, maybe it isn't open today?
                text = "Ups, algo ha salido mal:\n" + data["message"]
            else:
                text = (
                        "Hoy de primer plato hay:\n" +
                        '\n'.join(['- ' + p for p in data["primer"]]) +
                        "\n\nSegundo plato:\n" +
                        '\n'.join(['- ' + p for p in data["segundo"]])
                )
        except:
            text = (
                    "El servidor no está disponible ahora, vuelve a intentarlo más tarde.\n\nEl menú semanal se puede ver manualmente aquí:\n" +
                    "https://www.etsisi.upm.es/escuela/servicios/cafeteria"
            )

        bot.sendMessage(update.message.chat_id, text=text)


def busy_rooms_command(bot, update):
    if is_call_available("busy_rooms", update.message.chat_id, 180):
        log_message(update)
        logger.info("Getting busy rooms")
        options = webdriver.ChromeOptions()
        options.add_argument('--headless')
        options.add_argument('--window-size=1033,558')
        driver = webdriver.Chrome(executable_path='resources/chromedriver_linux', chrome_options=options)
        driver.get(settings.url_busy_rooms)
        driver.get_screenshot_as_file('resources/ocupacion_salas_biblioteca.png')
        driver.close()
        text = 'La ocupación de las salas de trabajo en grupo de la Biblioteca Universitaria Campus Sur es la ' \
               'siguiente: '
        bot.sendMessage(chat_id=update.message.chat.id, text=text, parse_mode=telegram.ParseMode.HTML)
        bot.send_photo(chat_id=update.message.chat_id, photo=open('resources/ocupacion_salas_biblioteca.png', 'rb'))


def schedule_command(bot, update, args):  # Add arguments for checking other's group schedule
    global chat_ids_list
    global schedule_list

    def schedule_parser(schedule):
        parsed_schedule = [""]
        schedule_keys = sorted(schedule, key=lambda s: int(s.split(":")[0]))
        for hour in schedule_keys:
            parsed_schedule.append("A las %sh -> %s" % (hour, schedule[hour]))
        return "\n".join(parsed_schedule)

    if is_call_available(update.message.text, update.message.chat_id, 60):
        log_message(update)
        try:
            group = ""
            day_index = datetime.datetime.today().weekday()
            if args:
                args = [str(x.upper()) for x in args]
                if re.match(r"G[TM][1-4]{2}", args[0]):  # First argument is a group identificator
                    group = args[0]
                    if len(args) > 1:  # If a second paramenter exists
                        if re.match(r"[LMXJV]", args[1]) or args[1].decode(
                                'utf-8') in codedays.keys():  # Two inputs: group and daycode
                            day_index = codedays[str(args[1].decode('utf-8'))]
                else:
                    if not re.match(r"[LMXJV]", args[0]) or not args[0].decode(
                            'utf-8') in codedays.keys():  # Two inputs: group and daycode
                        if update.message.chat_id < 0:
                            bot.send_message(chat_id=update.message.chat_id,
                                             text="Día de la semana inválido. Debes introducir martes/M, miércoles/X, jueves/J, viernes/V")
                            return
                        else:
                            bot.send_message(chat_id=update.message.chat_id,
                                             text="Debes especificar un grupo. <i>Por ejemplo: /horario gt11</i>",
                                             parse_mode=telegram.ParseMode.HTML)
                            return
                    day_index = codedays[str(args[0].decode('utf-8'))]

            if update.message.chat_id < 0:  # ID's below 0 are groups.
                group = update.message.chat.title.replace(" ETSISI", "")  # get group from chat title

            if group == u"GX3" or group == u"GX4":
                bot.send_message(chat_id=update.message.chat_id,
                                 text="Esta función no está definida para este grupo.")
            elif day_index in range(1, 5):
                text = schedule_parser(schedule_list[group][str(day_index)])
                text = "Horario del " + weekdays[day_index] + " para " + group + ":" + text
                bot.send_message(chat_id=update.message.chat_id, text=text)
            elif day_index == datetime.datetime.today().weekday():  # Check if user input is from 'today'
                bot.send_message(chat_id=update.message.chat_id,
                                 text="Hoy " + weekdays[day_index] + " no hay clases.")
            else:
                bot.send_message(chat_id=update.message.chat_id,
                                 text="El " + weekdays[day_index] + " no hay clases.")
        except:
            tb = traceback.format_exc()
            bot.send_message(chat_id=update.message.chat_id,
                             text="No he podido procesar tu solicitud de horario.\n\nERROR:\n" + str(
                                 tb) + "\n\nPor favor, reenvía este error a @nestoroa.")
    else:
        delete_message(bot, update)


def teacher_command(bot, update, args):  # Add arguments for checking other's group schedule
    global chat_ids_list
    global teacher_list

    def teacher_parser(teacher):
        parsed_teachers = [""]
        # schedule_keys = sorted(schedule, key=lambda s: int(s.split(":")[0]))
        for k, v in zip(teacher.keys(), teacher.values()):
            parsed_teachers.append("<b>%s</b> -> %s" % (k, v))
        return "\n".join(parsed_teachers)

    if is_call_available("teacher", update.message.chat_id, 180):
        log_message(update)
        try:
            group = ""
            if args:
                group = args[0].upper()
            if update.message.chat_id < 0:  # ID's below 0 are groups.
                group = update.message.chat.title.replace(" ETSISI", "")  # get group from chat title

            if group is not "":
                text = teacher_parser(teacher_list[group])
                text = "Lista de profesores para el primer semestre de " + group + ":\n" + text
            else:
                text = "No me has dicho a qué grupo te refieres."
            bot.send_message(chat_id=update.message.chat_id, text=text, parse_mode=telegram.ParseMode.HTML)

        except:
            tb = traceback.format_exc()
            bot.send_message(chat_id=update.message.chat_id,
                             text="No he podido procesar tu solicitud de profesores.\n\nERROR:\n" + str(
                                 tb) + "\n\nPor favor, reenvía este error a @nestoroa.")
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
        dispatcher.add_handler(CommandHandler('calendario', calendar_command))
        dispatcher.add_handler(CommandHandler('menu', menu_command))
        dispatcher.add_handler(CommandHandler('horario', schedule_command, pass_args=True))
        dispatcher.add_handler(CommandHandler('profesores', teacher_command, pass_args=True))
        dispatcher.add_handler(CommandHandler('salasbiblioteca', busy_rooms_command))
        dispatcher.add_handler(MessageHandler(Filters.status_update.new_chat_members, delete_message))
        dispatcher.add_handler(MessageHandler(Filters.status_update.new_chat_title, delete_message))
        dispatcher.add_handler(MessageHandler(Filters.status_update.left_chat_member, delete_message))
        dispatcher.add_error_handler(error_callback)

    except Exception as ex:
        logger.exception("Error al conectar con la API de Telegram.")
        quit()

    updater.start_polling()
    logger.info("ETSISI Bot: Estoy escuchando.")
    updater.idle()
