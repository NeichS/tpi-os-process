# tpi-os-process

## Guia para iniciar el programa

Dependencias minimas:

- go
- make (opcional, sirve para compilar y ejecutar a traves de un comando mas corto) ```sudo apt install make``` 

Para ejecutar el programa puede hacerlo simplemente ejecutando el makefile:

    make run

Si no desea instalar make:

    go build -o bin/main cmd/main/main.go


> Compilará y ejecutará el programa automaticamente con los registros que existen por defecto.

Si desea ingresar un archivo debe colocarlo dentro de la carpeta csv-files y luego ejecutar el comando de make:

    make run PARAM=<nombre-del-archivo.csv>
    ./bin/main <nombre-del-archivo.csv> #si no utiliza make

El archivo generado con los logs del programa se encontrará dentro de la carpeta **output**

## Supuestos:

**El tiempo avanza por vuelta, con estome refiero a la ocurrencia secuencial de las siguientes operaciones:**  

1. Corriendo a Terminado.
2. Corriendo a Bloqueado.
3. Corriendo a Listo.
4. Bloqueado a Listo.
5. Nuevo a Listo.
6. Finalmente el despacho de Listo a Corriendo.
 
- bloqueado a listo ocurre instantaneamente.

- En cada uno de los eventos si no se hace nada no se suma tiempo .

- Los procesos ejecutan sus rafagas de cpu y luego las rafagas de I/O.

- Cuando un proceso termine su ultima rafaga de I/O este debera ser despachado una vez mas para pasar al estado terminado.

- El tip es un tiempo que se agrega si y solo si el sistema operativo a aceptado un proceso de nuevo a listo.

- En srt cuando un proceso ejecuta su ultima operacion de entrada y salida, va a volver a requerir procesador pero su remainging time va a ser 0.

- El tiempo de retorno se empieza a tomar desde el tiempo de arribo del proceso. Tiempo en el que termina - Tiempo de arribo.

- Round robin, si TCP >= quantum, se produce un bucle infinito ya que ningun proceso llega a poder emitir una rafaga. 
