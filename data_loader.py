#!/usr/bin/env python
# -*- coding: utf-8 -*-
from logger import get_logger
import json

logger = get_logger("data_loader")


class DataLoader:

    def __init__(self):
        global data_and_settings
        try:
            private_data = json.load(open('private-data.json'), encoding="utf-8")
            data = json.load(open('data.json'), encoding="utf-8")
            data_and_settings = private_data.copy()
            data_and_settings.update(data)
        except:
            logger.exception("Error al cargar el JSON de configuración")
        else:
            logger.info("JSON cargado con éxito")
            self.telegram_token = data_and_settings["telegram_token"]
            self.help_string = data_and_settings["strings"]["help"]
            self.github_url = data_and_settings["strings"]["github_url"]
            self.admin_password = data_and_settings["admin_password"]
            self.degrees = data_and_settings["degrees"]
            self.etsisi_urls = data_and_settings["etsisi_urls"]
            self.upm_jsons = data_and_settings["upm_jsons"]
