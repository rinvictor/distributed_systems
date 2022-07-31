**Víctor Rincón Yepes.**
Sistemas Distribuidos.
# Práctica 3. Los almacenes de Santa Claus.
Se trata de una modificación de la práctica anterior en la que se incluyen dos cachés de duendes, un almacén principal de regalos con 90 regalos repartidos en 3 filas.

Tal y como se especifica en la práctica se añaden una serie de restricciones específicas que se indican en el enunciado de la práctica.

## Solución propuesta:
Los líderes tendrán acceso tanto al almacén como a las cachés para que puedan se actualizadas. Estas estructuras se pasan por referencia a los líderes y son inicializadas en el programa principal.
Para que cada líder conozca que caché le corresponde se le pasa el valor "tag".

En el almacén se representa si hay un regalo o no en esa posición, al principio estará completamente lleno. En cambio, en las cachés se representa si esa posición está disponible.
Este es el motivo por el que los valores del almacén y las cachés están invertidos.

Aquí podemos ver un ejemplo:
```
--ESTADO DEL ALMACÉN Y LAS CACHÉS--
0000000000X0X0000000XXX000X0X0
0XX0000X000000XX00X000000000X0
000000X000XX00X000000XX0000000

XXXXXXXXXX0X0XXXXXXX000XXX0X0X
X00XXXX0XXXXXX00XX0XXXXXXXXX0X
XXXXXX0XXX00XX0XXXXXX00XXXXXXX

XXXXXXXXXX0X0XXXXXXX000XXX0X0X
X00XXXX0XXXXXX00XX0XXXXXXXXX0X
XXXXXX0XXX00XX0XXXXXX00XXXXXXX
```
Implementación de la caché y la fábrica de juguetes:
>var toyStore [Rows][Spaces]bool //Almacén de juguetes

>var cache [NBattalions][Rows][Spaces]bool //Una caché por batallon

Se han escrito las siguientes funciones para dar solución al problema propuesto:

>func santa(reindeersChan chan int, elvesChan chan int, leadersChan chan int).

Como vemos, recibe tres canales, reindeersChan, elvesChan, leadersChan.

Los renos le despiertan por reindeersChan y los elfos lo harán por elvesChan. El tercer canal se usa para que los líderes avisen de que han terminado de preparar todos los regalos.

Si son los renos quienes despiertan a Santa se hace alert = 1, así los elfos sabrán que el programa termina. Esta variable compartida es necesaria ya que de no usarse los elfos quedarían intentado avisar a Santa para que los ayude y este no respondería.

Los elfos leerán esta variable cuando vayan a pedir ayuda a Santa. Cuando sean los líderes quienes la lean avisarán a todos los elfos:

```
for i := 0; i < D; i++ {
	elves <- 0 //indicamos a los elfos que hemos terminado
}
```



>func elf(toysChan chan int, elvesChan chan int)

Tiene problemas en uno de cada tres juguetes, esto se consigue usando la función rand.Intn(),que recibe un parámetro entero y genera un número aleatorio comprendido entre cero y este.
A diferencia de la práctica dos, en el caso de que un elfo tenga problemas con un regalos podemos distinguir dos casos:

 - Es el primer o el segundo elfo. No se quedan esperando a que Santa los ayude, sino que dejan ese regalo y continuan con otro diferente.

 - Es el tercer elfo. Espera a que Santa termine con su ayuda y y actualiza el contador de elfos con problemas y regalos terminados, incrementándolo en tres.

Como hemos comentado antes para evitar un bloqueo de los elfos se realiza la compobación de la variable compartida "alert".

La gorutina termina cuando lo comunica el líder:
```
if (<-toysChan) == 0{ //el lider dice que no quedan juguetes, terminamos
	break
}
```

>func reindeer(reindeersChan chan int)

Comprueba si es el último reno, en este caso despierta a Santa, usando el canal que recibe.

Se espera hasta que Santa termine de hacer sus cosas (wait.Lock())

>func leader(tag int, caches *[NBattalions][Rows][Spaces]bool, toyStore *[Rows][Spaces bool, elvesChan chan int, leadersChan chan int)

Recibe "tag" para poder conocer cual es su caché. Recibe la fábrica de juguetes y los canales que usará.

En primer lugar la función entra en un bucle del que saldrá cuando alert==1, es decir, que Santa ya no está, o cuando se hayan terminado de fabricar todos los juguetes.
A continuación, se realiza la selección de filas, para ello se usa una variable en la que se anotan las filas ya seleccionadas.

Tras elegir una fila, se elige al azar uno de los treinta regalos. Accede a su caché para comprobar su disponibilidad. Si está disponible, lo coge. Si no está disponible, prueba con el de índice siguiente.

Este es el código usado para esa parte del programa:

```
if !toys[row] { //Si esa fila no se ha elegido aun
	toys[row] = true //La marcamos como elegida
	n := rand.Intn(Spaces) //Un regalo al azar
	for j := 0; j < Spaces; j++{ //Hasta que un regalo que aún esté disponible se elija 
		pos := (n + j) % Spaces //Paso al siguiente cuando esté vacío
		if !(*caches)[tag][row][pos] { //Compruebo en mi caché si el regalo está disponible
			toyStore[row][pos] = false //Lo retiro del almacén
			//Actualizo ambas cachés
			for k := 0; k < NBattalions; k++ {
				(*caches)[k][row][pos] = true
			}
			fmt.Println("LíderB" + strconv.Itoa(tag) + ": ¡Va un regalo!")
			elvesChan <- 1 //Le mando un regalo a los elfos
			break
		}else if j == (Spaces-1){
			full++ //La fila está llena
		}
	}
```

Se actualizan las cachés y se envía un regalo a los elfos.

Una vez se han entregado todos los regalos o Santa se ha ido, informa a los elfos de su batallón para que dejen de trabajar. En el primer caso además informará a Santa y a la función que crea los renos para que terminen también (alert=2).

Hay dos funciones más que se usan para lanzar las gorutinas con las funciones anteriormente descritas, estas son muy similares a las de la praćtica anterior pero tienen alguna modificación como podemos ver a continuación:

>func launchElves(toyStore *[Rows][Spaces]bool, caches *[NBattalions][Rows][Spaces]bool, elvesChan chan int, leadersChan chan int)

Crea los batallones de elfos y los líderes de cada batallón.

>func launchReindeers(reindeersChan chan int)

Lanza todos los renos pero antes comprueba si alert==2, es decir, si se ha terminado de fabricar todos los regalos, ya que en este caso debería morir.


Además, se espera la finalización de las gorutinas para acabar la ejecución. Para ello:

>wg.Done()

En el main se lanzan las diferentes gorutinas que hemos explicado anteriormente y se imprime cuando conviene el estado de las cachés y del almacén de juguetes usando la función:

>printStatus(input *[Rows][Spaces]bool)
## Esquemas.

>Esquema de la función santa:

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/pfinal/fsanta_pfinal.png" width="280" height="280" alt="santa">

>Esquemas de la función líder:

Nota: Para que sea más comprensible gráficamente se ha dividido en tres casos que provendrían del estado del programa que se indica en cada gráfico.

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/pfinal/santa_caso1_pfinal.png" width="280" height="280" alt="elfo">

Caso 1.

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/pfinal/santa_caso2_pfinal.png" width="280" height="280" alt="elfo">

Caso 2.


<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/pfinal/santa_caso3_pfinal.png" width="280" height="280" alt="elfo">

Caso 3.

>Esquema de la función elfo:

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/pfinal/elfo_pfinal.png" width="280" height="280" alt="elfo">

>Esquema de la función reno:

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/pfinal/reno_pfinal.png" width="280" height="280" alt="elfo">

## Ejemplo de ejecución.
```
--ESTADO DEL ALMACÉN--
XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

Santa duerme...
LíderB1: ¡Va un regalo!
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 0 terminado
Juguete 1 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
un reno llega, van 1
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 2 terminado
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: Tengo problemas con este juguete
Juguete 3 terminado
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: Tengo problemas con este juguete
Elfo: Tengo problemas con este juguete
Elfo: A despertar al jefe para que nos ayude
Santa: Voy a ayudar a estos duendecillos
Juguete 4 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: Tengo problemas con este juguete
Santa duerme...
Elfo: Ahí van 3 juguetes más.
un reno llega, van 2
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 8 terminado
Juguete 9 terminado
Juguete 10 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 11 terminado
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 12 terminado
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
un reno llega, van 3
Juguete 13 terminado
Juguete 14 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: Tengo problemas con este juguete
Juguete 15 terminado
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 16 terminado
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
un reno llega, van 4
Juguete 17 terminado
Elfo: Tengo problemas con este juguete
Elfo: A despertar al jefe para que nos ayude
Santa: Voy a ayudar a estos duendecillos
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 18 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Santa duerme...
Elfo: Ahí van 3 juguetes más.
Juguete 22 terminado
Elfo: Tengo problemas con este juguete
Elfo: Tengo problemas con este juguete
Elfo: Tengo problemas con este juguete
Elfo: A despertar al jefe para que nos ayude
Santa: Voy a ayudar a estos duendecillos
Juguete 23 terminado
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
un reno llega, van 5
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: Tengo problemas con este juguete
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Santa duerme...
Elfo: Ahí van 3 juguetes más.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 27 terminado
Juguete 28 terminado
Juguete 29 terminado
Juguete 29 terminado
un reno llega, van 6
Juguete 31 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 32 terminado
Juguete 33 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: Tengo problemas con este juguete
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB0: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 34 terminado
un reno llega, van 7
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Juguete 35 terminado
Juguete 36 terminado
Juguete 37 terminado
Elfo: Tengo problemas con este juguete
Elfo: A despertar al jefe para que nos ayude
Santa: Voy a ayudar a estos duendecillos
Juguete 38 terminado
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
un reno llega, van 8
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Santa duerme...
Elfo: Ahí van 3 juguetes más.
Juguete 42 terminado
Elfo: Tengo problemas con este juguete
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Elfo: Tengo problemas con este juguete
LíderB1: ¡Va un regalo!
Elfo: ¡Ya tengo trabajo!, recibí un juguete.
Rudolf: Soy el reno que faltaba, perdón por llegar tarde...¡A repartir!
Santa: Ya están todos los renos aquí, ¡creía que no llegarían nunca!
Santa: Aviso a los elfos de que me voy (alert=1)
Santa duerme...
Juguete 43 terminado
LíderB0: Santa se ha ido, terminad lo que estéis haciendo y a casa
LíderB1: Santa se ha ido, terminad lo que estéis haciendo y a casa
Elfo: Tengo problemas con este juguete
Elfo: A despertar al jefe para que nos ayude
Elfo: Santa ya se se ha ido...
Juguete 44 terminado
Fin. ¡Feliz Navidad!
--ESTADO DEL ALMACÉN Y LAS CACHÉS--
X0X0XX00XXXX0XX000XX0000000XXX
X00000000XX00XXX000XXXX000X0XX
X0XX00000000XXX00X0000XXXXX0XX

0X0X00XX0000X00XXX00XXXXXXX000
0XXXXXXXX00XX000XXX0000XXX0X00
0X00XXXXXXXX000XX0XXXX00000X00

0X0X00XX0000X00XXX00XXXXXXX000
0XXXXXXXX00XX000XXX0000XXX0X00
0X00XXXXXXXX000XX0XXXX00000X00
```



