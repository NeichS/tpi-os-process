# tpi-os-process


## Supuestos:

- La unidad de tiempo sera en milisegundos

- Tiempo de retorno normalizado: ```tiempo de retorno / tiempo de servicio``` (**a tener en cuenta para el tpi**)

- El valor de la prioridad externa [0 , 5] siendo 0 el mas prioritario

- El tiempo avanza por vuelta

**Con vuelta me refiero a la ocurrencia de todos los eventos siguientes:**  

1. Corriendo a Terminado.
2. Corriendo a Bloqueado.
3. Corriendo a Listo.
4. Bloqueado a Listo.
5. Nuevo a Listo.
6. Finalmente el despacho de Listo a Corriendo.


- bloqueado a listo ocurre instantaneamente

- En cada uno de los eventos si no se hace nada no se suma tiempo 

- El tip es un tiempo que se agrega si y solo si el sistema operativo a aceptado un proceso de nuevo a listo