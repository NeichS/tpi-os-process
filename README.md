# tpi-os-process

## Guia para iniciar el programa

Dependencias minimas:

- go
- make (opcional, sirve para compilar y ejecutar a traves de un comando mas corto) ```sudo apt install make``` 

Crear una carpeta 

Para ejecutar el programa puede hacerlo simplemente ejecutando el makefile:

    make run

Si no desea instalar make:

    go build -o bin/main cmd/main/main.go


> Compilará y ejecutará el programa automaticamente con los registros que existen por defecto.

Si desea ingresar un archivo debe colocarlo dentro de la carpeta csv-files y luego ejecutar el comando de make:

    make run PARAM=<nombre-del-archivo.csv>
    ./bin/main <nombre-del-archivo.csv> #si no utiliza make

El archivo generado con los logs del programa se generará en raiz con el nombre de logs.txt 

## Supuestos:
 
- bloqueado a listo ocurre instantaneamente.

- En cada uno de los eventos si no se hace nada no se suma tiempo .

- Los procesos ejecutan sus rafagas de cpu y luego las rafagas de I/O.

- Cuando un proceso termine su ultima rafaga de I/O este debera ser despachado una vez mas para pasar al estado terminado donde empezara a ejecutar su TFP. 

- Un proceso inicia la ejecucion de su TIP cuando pasa de estado nuevo a listo.

- En srt cuando un proceso ejecuta su ultima operacion de entrada y salida, va a volver a requerir procesador pero su remainging time va a ser 0.

- El tiempo de retorno se empieza a tomar desde el tiempo de arribo del proceso. Tiempo en el que termina - Tiempo de arribo.

- Un proceso para pasar a estado terminado primero debe volver a estar en estado running, y como en mi programa yo considero que una rafaga esta compuesta por las rafagas de cpu necesarias + las rafagas de entrada y salida necesarias, entonces un proceso para terminar debera volver a ser despachado para inciar su TFP o termina instantaneamente si tfp = 0