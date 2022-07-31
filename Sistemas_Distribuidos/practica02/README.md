**Víctor Rincón Yepes.**
Sistemas Distribuidos.
# Práctica 2. Los almacenes de Santa Claus.
Una modificación de la práctica anterior en la que se incluyen dos cachés de duendes, un almacén principal de regalos con 90 regalos repartidos en 3 filas.

Tal y como se especifica en la práctica se añaden una serie de restricciones específicas que se indican en el enunciado de la práctica.

## Solución propuesta:

Es un problema de sincronización que tendrá su final en el momento en el que llega el último reno y comienza el reparto.

Como mecanismo de sincronización he optado por utilizar mutex y variables compartidas entre las distintas posibilidades que había.

Se han escrito tres funciones para representar el comportamiento de cada uno de los personajes implicados en el problema:

>func santa(R int), recibe como parámetro el número de renos. Comportamiento de Santa Claus.

En el enunciado se especifica que Santa debe dar prioridad a la llegada de los renos y que en el caso de que todos los renos esteń listos y que haya elfos que tengan problemas para terminar su juguete simplemente el trineo saldrá y los elfos no recibirán ayuda.

Esto se modela usando un if y primero comprobando que quien le ha llamado sea elúltimo de los renos y haciendo un break:
```
sledge.Lock()
//alguien le despierta, los renos tienen prioridad
if nreindeers == R {
    fmt.Println("Santa: Ya están todos los renos aquí, ¡creía que no llegarían nunca!")
    break
}else{
    fmt.Println("Santa: Voy a ayudar a estos duendecillos")
    time.Sleep(time.Duration(rand.Float64()*(maxtime-mintime)+mintime)* time.Second) //aleatorio entre 2 y 5 seg
    elfWait.Unlock() //termino de ayudarlos
}
sledge.Unlock()
```

>func elf(), comportamiento de un elfo.

Tiene problemas en uno de cada tres juguetes, esto se consigue usando la función rand.Intn(),que recibe un parámetro entero y genera un número aleatorio comprendido entre cero y este, hay un ejemplo del comportamiento de esta función más adelante, explicando sus mutex.

>func reindeer(R int), recibe como parámetro el número de renos.

Comprueba si es el último reno, en este caso hace santasBed.Unlock() para despertar a Santa.

Hay dos funciones más que se usan para lanzar las gorutinas con las funciones anteriormente descritas:
>func launchElves(D int), lanza elfos, todos seguidos.

>func launchReindeers(R int), lanza renos con un intervalo mínimo de 5 segundos.

En el main se indican el número de elfos y renos, se inicializan los mutex y se llama a Santa y a las funciones que lanzan a elfos y renos.

Además, se espera la finalización de las gorutinas para acabar la ejecución.
```
func main() {
	D := 12 //número de elfos
	R := 9 //número de renos
	wg.Add(D+R+1) //número de gorutinas
	rand.Seed(time.Now().UTC().UnixNano()) //aleatoreidad, necesario, si no no serían 
											//los mismos números siempre en las funciones rand

	//init
	santasBed.Lock()
	elfWait.Lock()

	go santa(R)
	go launchReindeers(R)
	go launchElves(D)

	wg.Wait() //Espero que acaben todas las gorutinas lanzadas
	fmt.Println("Fin. ¡Feliz Navidad!")
}
```

Para solucionar este problema he usado los mutex:
>*santasBed*: Representa que Santa esta descansando

>*elfWait*: Los Elfos tienen que esperar a que Santa los ayude

>*factory*: Lo uso a la hora de comprobar si podemos o no llamar a Santa (paramos la producción)

>*sledge*: El trineo controla la variable compartida renos.Cuando los renos se "enganchan" al trineo

Las variables compartidas utilizadas son:
>*nelvesintrouble*: Representa el número de elfos que están encontrando problemas para acabar un regalo.

>*nreindeers*: representa el número de renos que ya ha llegado.

Notar que la variable compartida **_nelvesintrouble_** está controlada por el mutex **_factory_** y **_nreindeers_** controlada por **_sledge_**.

**_santasBed_** como hemos indicado anteriormente es un mutex que sirve para despertar (Unlock) o hacer que duerma (Lock) Santa Claus.

Adicionalmente, se usa el mutex **_elfwait_** para bloquear al elfo (Lock) que llama a Santa mientras que este esté ayudando al grupo, cuando termine Santa, hará elfWait.Unlock(), este elfo pone la variable compartida a 0 y continúa su ejecución. Esto se hace ya que no se debe permitir el acceso a **_nelvesintrouble_** hasta que este reseteada.

Lo podemos observar aquí:
```
		fmt.Println("Elfo: Tengo problemas con este juguete")
		factory.Lock() //paro para comprobar
		nelvesintrouble++
		if nelvesintrouble == minelvesintrouble{
			fmt.Println("Elfo: A despertar al jefe para que nos ayude")
			santasBed.Unlock() //Despierto
			elfWait.Lock() //esperan a que termine de ayudar a estos
			//reseteo el conteo
			nelvesintrouble = 0
		}
		factory.Unlock()
```

Para finalizar correctamente espero a la finalización de las gorutinas lanzadas usando:

>wg.Wait(), variable global wait group

En el main se indica el número de rutinas por las que debe esperar y en cada finalización de una de estas se hace wg.Done().

Cuando todas (renos, elfos y santa) han terminado el main finaliza su ejecución.


## Esquemas.

>Esquema funcionamiento general:

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/practica02/general2.png" width="250" height="250" alt="general">

>Esquema de la función santa:

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/practica02/santa2.png" width="280" height="280" alt="santa">

>Esquema de la función elf:

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/practica02/elfo2.png" width="280" height="280" alt="elfo">

>Esquema de la función reno:

<img src="/home/vrincon/Dropbox/Sistemas_Distribuidos/practica02/reno2.png" width="280" height="280" alt="reno">




## Ejemplo de ejecución.
```
$>go run practica2.go 
Santa duerme...
juguete terminado
Elfo: Tengo problemas con este juguete
un reno llega, van 1
juguete terminado
un reno llega, van 2
Elfo: Tengo problemas con este juguete
un reno llega, van 3
un reno llega, van 4
juguete terminado
juguete terminado
un reno llega, van 5
juguete terminado
juguete terminado
juguete terminado
un reno llega, van 6
Elfo: Tengo problemas con este juguete
Elfo: A despertar al jefe para que nos ayude
Santa: ¡Alguien me ha despertado!
Santa: Voy a ayudar a estos duendecillos
Santa duerme...
un reno llega, van 7
Elfo: Tengo problemas con este juguete
juguete terminado
un reno llega, van 8
Rudolf: Soy el reno que faltaba, perdón por llegar tarde...¡A repartir!
Santa: ¡Alguien me ha despertado!
Santa: Ya están todos los renos aquí, ¡creía que no llegarían nunca!
Fin. ¡Feliz Navidad!
```



