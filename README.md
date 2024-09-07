# tpi-os-process


## Supuestos:

- La unidad de tiempo sera en milisegundos

- Tiempo de retorno normalizado: ```tiempo de retorno / tiempo de servicio``` (**a tener en cuenta para el tpi**)

- El valor de la prioridad externa [0 , 5] siendo 0 el mas prioritario

- El tiempo avanza por cada uno de los eventos que efectivamente suceden

**Con evento me refiero a los listados en la consigna:**  

1. Corriendo a Terminado.
2. Corriendo a Bloqueado.
3. Corriendo a Listo.
4. Bloqueado a Listo.
5. Nuevo a Listo.
6. Finalmente el despacho de Listo a Corriendo.


- bloqueado a listo ocurre instantaneamente

- En cada uno de los eventos si no se hace nada no se suma tiempo 

- El tip es el tiempo en el que el SO se fija si hay procesos nuevos para agregarlos a listo