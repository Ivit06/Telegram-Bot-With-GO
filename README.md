# FujiHat

Aquest tutorial et guiarà a través dels passos necessaris per a configurar i executar el nostre bot. Aquest bot està programat en Go.

Una de les habilitats més importants dins del món de la informàtica és la capacitat d'automatitzar tasques. Això ens permet estalviar tant temps com recursos. 

Davant d’aquesta necessitat i com a administradors de sistemes hem vist adequat realitzar un projecte on oferim una eina automatitzada, útil i senzilla.

Les seves funcions principals són les següents:

1.  Monitoratge d'instàncies.
2.  Autodescobriment d'instàncies.
3.  Administració dels usuaris autoritzats.
4.  Executa consultes a mètriques.
5.  Enviar alertes relacionades amb els components de les instàncies.

## Instalación

1. **Go:**
   
    ```bash
     apt install golang-go
    ```
   
2.  **Descarga las dependencias:**


    Les dependències ja estan incloses en el repositori.
    Si us dona algun tipus de problema es poden actualitzar les dependències amb el següent comando:

    ```bash
     go mod tidy
    ```
    
3.  **Ngrok:**

4.  **Base de dades (master/slave):**

       Es opcional realizar la configuración de las bases de datos en una única instancia donde no haya redundancia de datos. Para ello se deberán de hacer modificaciones en el código para que únicamente se traten los datos desde la misma instancia.

5.  **Clonar Repositori:**
    ```bash
    git clone [https://github.com/sindresorhus/del](https://github.com/sindresorhus/del)
    cd [Nombre de tu Proyecto]
    ```

    Si no tienes Git instalado, puedes descargar el código fuente como un archivo ZIP y extraerlo.


6. **Prometheus**


## Configuración

1.  **Crea un archivo de configuración:**

    Copia el archivo `config.example.yaml` a `config.yaml` y edítalo con tus propios valores de configuración.

    ```bash
    cp config.example.yaml config.yaml
    nano config.yaml
    ```

    (o usa tu editor de texto preferido).

2.  **Configura las variables de entorno (opcional):**

    Algunas variables de configuración pueden establecerse a través de variables de entorno. Consulta el archivo `config.example.yaml` para obtener más detalles.

## Ejecución

1.  **Compila y ejecuta la aplicación:**

    ```bash
    go run main.go
    ```

    O, para compilar un binario ejecutable:

    ```bash
    go build -o [nombre del ejecutable]
    ./[nombre del ejecutable]
    ```

2.  **Accede a la aplicación:**

    [Describe cómo acceder a la aplicación, por ejemplo, a través de un navegador web o una interfaz de línea de comandos].

## Uso

[Describe cómo usar la aplicación, incluyendo ejemplos de comandos o capturas de pantalla].

## Ejemplos

[Proporciona ejemplos de casos de uso comunes].

## Contribución

Si deseas contribuir a este proyecto, sigue estos pasos:

1.  Haz un fork del repositorio.
2.  Crea una rama para tu contribución: `git checkout -b mi-contribucion`.
3.  Realiza tus cambios y haz commit: `git commit -m "Descripción de los cambios"`.
4.  Sube tus cambios al repositorio remoto: `git push origin mi-contribucion`.
5.  Crea un pull request.

## Licencia

Este proyecto está bajo la licencia [Nombre de la licencia]. Consulta el archivo `LICENSE` para obtener más detalles.
