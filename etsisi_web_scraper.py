#!/usr/bin/env python
# -*- coding: utf-8 -*-
import io
import json
from bs4 import BeautifulSoup
import requests


# Adquiere enlaces a las guías docentes alojadas en la web de la escuela
def degree_json_scraper():
    def get_degree_guides(degree_code):
        print("Getting guides for " + str(degree_code))
        requested_url = data['urls'][degree_code]
        html = requests.get(requested_url)
        soup = BeautifulSoup(html.text, 'html.parser')
        asignaturas_dict = {
        }
        for idx, row in enumerate(soup.find_all('tr')):
            contents = row.contents
            codigo_asignatura = str(contents[1].string).strip("\n").strip("\t")
            if codigo_asignatura == 'None':
                continue
            # Iniciamos el diccionario
            asignaturas_dict[codigo_asignatura] = {}
            # Guardamos nombre
            asignaturas_dict[codigo_asignatura]["nombre"] = str(contents[3].string).strip("\n").strip("\t")
            # Guardamos enlaces
            links = ""
            for link_count, link in enumerate(contents[5].findAll("a")):
                if link_count == 0:
                    links = "https://etsisi.upm.es" + str(link.get("href"))
                else:
                    links = "Enlace 1: " + links + " \n Enlace 2: https://etsisi.upm.es" + str(link.get("href"))
                asignaturas_dict[codigo_asignatura]["guia_docente"] = links
        return asignaturas_dict

    with io.open('data.json', 'r', encoding='utf8') as data_file:
        data = json.load(data_file)
    degrees_guides = {}
    for degree in data["degrees"]:
        degrees_guides[degree] = get_degree_guides(degree)

    # """We could save this in a JSON"""
    # with io.open("data_dumped.json", 'w', encoding='utf8')  as data_file:
    #     json.dump(degrees_guides, data_file, ensure_ascii=False)

    return degrees_guides


# Listado de noticias disponible en la web de la UPM
def news_json_scraper():
    html = requests.get("https://etsisi.upm.es/noticias")
    soup = BeautifulSoup(html.text, 'html.parser')
    news_dict = {}
    for idx, row in enumerate(soup.find(id='main-content').findAll('a')):
        news_dict[idx] = {}
        news_dict[idx]["text"] = row.string
        news_dict[idx]["a-link"] = '<a href="https://etsisi.upm.es' + row.get("href") + '">' + row.string + '</a>'
        news_dict[idx]["link"] = 'https://etsisi.upm.es' + row.get("href")
    return news_dict


# Listado de eventos disponible en la web de la upm
def events_json_scraper():
    html = requests.get("https://etsisi.upm.es")
    soup = BeautifulSoup(html.text, 'html.parser')
    events_dict = {
    }
    for idx, row in enumerate(soup.find(id='block-views-calendario-block-2--2').findAll('a')):
        events_dict[idx] = {}
        events_dict[idx]["text"] = row.string
        events_dict[idx]["a-link"] = '<a href="https://etsisi.upm.es' + row.get("href") + '">' + row.string + '</a>'
        events_dict[idx]["link"] = 'https://etsisi.upm.es' + row.get("href")
    return events_dict


def avisos_json_scraper():
    html = requests.get("https://etsisi.upm.es/alumnos/avisos")
    soup = BeautifulSoup(html.text, 'html.parser')
    avisos_dict = {}

    for idx, row in enumerate(soup.find(id='main-content').findAll('a')):
        if idx > 5:
            continue
        avisos_dict[idx] = {}
        avisos_dict[idx]["text"] = row.string
        avisos_dict[idx]["a-link"] = '<a href="https://etsisi.upm.es' + row.get("href") + '">' + row.string + '</a>'
        avisos_dict[idx]["link"] = 'https://etsisi.upm.es' + row.get("href")

    avisos_dict[6] = {}
    avisos_dict[6]["text"] = "Más avisos..."
    avisos_dict[6]["a-link"] = '<a href="https://etsisi.upm.es/alumnos/avisos">Más avisos...</a>'
    avisos_dict[6]["link"] = 'https://etsisi.upm.es/alumnos/avisos'
    return avisos_dict


if __name__ == "__main__":
    print(news_json_scraper())
    print(events_json_scraper())
    print(avisos_json_scraper())
