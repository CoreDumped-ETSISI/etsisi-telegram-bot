# etsisi-telegram-bot

## Información

Bot y alojado por la asociación universitaria [>Core Dumped_](https://coredumped.es)

## Contribuir

Para añadir comandos al bot, se crea una carpeta por "feature" en la carpeta `cmds`, y luego se enrutan los comandos en `router.go`.

El bot llama a microservicios alojados en la repo https://github.com/CoreDumped-ETSISI/uni-services

## Despliegue

Para desplegar el bot con todos los microservicios se ha proporcionado un `docker-compose.yml` para su facilidad. Simplemente cambie los valores `<replaceme>` por sus valores reales y ejecute `docker-compose up` en la carpeta del proyecto.

Para el desarrollo, solo hace falta cambiar el valor de la clave de API de telegram.

# [Colaboradores](https://github.com/CoreDumped-ETSISI/etsisi-telegram-bot/graphs/contributors)
