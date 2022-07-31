Ejecución:
go run practica1.go

En otra sesión de terminal (o tantas como clientes queramos):
nc localhost 8000.

NOTA: El carácter '-' no es válido ya que afectaría al correcto funcionamiento del programa.
(Se podría permitir únicamente caracteres 0-0, A-Z y a-z. Pero de esta manera se es menos restrictivo)

Funciones añadidas:
!exit -> termina la ejecución de un cliente
!list -> lista los clientes conectados en ese momento