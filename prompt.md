Quiero que actues como un senior fullstack developer. Experto en arquitectura y clean code.

Aunque me preguntes y te responda en español, todo el codigo generado, ejemplos, guias y tests, implemntalo todo en ingles.

Quiero que hagas una aplicacion mmultiplataforma hecha en ionic, es una aplicacion donde guardare un stock de productos, una despensa, para controlar el stock.

* Funcionalidades

- Quiero una primera version de una APP web (mbbile para despues) que sea responsive para verla desde el mobil, pero siendo web.
- Quiero que sea una aplicacion sencilla, un interfaz minimalista, con colores suaves y no muy brillantes.
- La web no tendra autenticacion, cualquier usuario anonimo puede acceder a todas sus fucnionalidad sin registro ni usuarios.
- La panatalla principal se vera el titulo de la app, pueude que alguna imagen para que no sea tan aburrido, y una pequeña descripcion, que dira que peudes hacer. Habra una textedit en el que podras añadir un titulo valido, el titulo valido sera el texto que el usuario desee, solo que filtrando caracteres raros o no permitivos dentro de una URL. Pues este titulo se usara para crear un enpoint por ejemplo si pongo "familia-suarez" se creara un endpoint web/despensas/familia-suarez. y abra un boton de crear. que creara este recurso y se dirigara hacia el nuevo endpoint.
- La seguna pantalla sera la pantalla de la despensa. esta pantalla tendra una lista de elementos predefinidos dentro de la app. cada elemento sera una imagen, por ejemplo una imagne de tomates, un texto "tomates", y un estado, que por ahora puede tener 3 estados, agotado color rojo, poco color amarillo y sufcienten color verde.
-Puedes seleccionar cada uno de los elementos, cada elemento, puedes ponerle categorias o tags, que serviran para filtrar, por ejemplo, categoria esencial, verdura, carnes, limpieza. las categorias tambien van a ser predefinidas.
- en la pantalla de la despensa podras ver los elementos segun el filtro, asi que todos los elementos empezaran sin ninguna categoria o filtro y luego segun el filtro o lista o vista, que quiera ver se organizaran o no
- habra un panel lateral donde se podran seleccion las vistas, por ahora habran 3 vistas, Primarios, secundarios, y Otros. Cada producto se podra poner en una de estas vistas. ( configurable en la configuracion del elemento)
- dentro de cada una de estas vistas se podra filtrar por categoria, asi se podra filtrar, cada filtro deberia poder escogerse rapidamente, no como un combobox, mas bien como una lista de iconos facilmente seleccionable. Ayudame a escoger las categorias, son productos de despensa, ayudame a  agruparlos entre 3-5 categorias, tmapoco quiero muchas.
- finalmente, en cada una de las vistas, habra un boton para generar un texto plano que se pueda copiar. en el que se generara una lista de todos los productos que no esten en stock. hbra un ehck box por si tmb quiero añdir a la lista elementos que tienen poco stock. la lista en texto plano la puedes dividir por categorias.

* Tecnologias

- QUiero que hagas esta app en un unico repositorio, osea un monolito backend, frontend e indra todo en un solo repo .
- toda la infra estructura quiero que te bases en /home/jsuarez/workspace/ductifact/backend y en /home/jsuarez/workspace/ductifact/infra para seguir la misma metodologia de desarrollo y delivery. Quiero sorbe todo que sean containeers de docker, que haya un makefile para poder montar la app y correrlo todo. Es importante que sigas mas o menos la misma infra para mas adelanto yo poder entender y hacer cambios necesarios.
- el framework que quiero que uses es IONIC, pero como soy nuevo en esto, no estoy sewguro que vas a necesitar de mas para desarrollarlo todo, en este caso analizalo y preguntame primero que tecnologias me recomiendas usar, estoy abierto a todo.

- es una app pequeña y local, no va a ser publica ni comercial, por lo que no quiero que hagas sobreeingenieria, solo la minima, es decir, si quieres hacer test, has los minimos no exageres. Problemas de seguridad los minimos, no hace falta bliondrlo a prueba de bombas

- en general el obbjetivo de este projecto es hacer una aplicacion muy sencilla, enfocada a prender un poco ionic, asi que si todo se puede hacer en ionic, mejor.

- Por ultimo trndras que hacer una pequeña guia de lo que has implementado, para entder lo que has hecho sobre todo sobre ionic y forntend pues donde mas quiero aporender a la vez que hago una aplicacion funcional y practica.

Analisa lee el codigo y planea lo que vas a implementar. Todas las preguntas y dudas que  tengas preguntamelas para discutirlas, y genera un plan si quieres en un .MD o como tu creas mejor, y no comienzes la implementacion hasta que yo te confirme todo.