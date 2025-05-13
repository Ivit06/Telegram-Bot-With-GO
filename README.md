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

## Instal·lació i Configuració

1. **Go:**
   
    ```bash
     apt install golang-go
    ```

   
2.  **Descàrrega les dependències:**

    Les dependències ja estan incloses en el repositori.
    Si us dona algun tipus de problema es poden actualitzar les dependències amb el següent comando:

    ```bash
     go mod tidy
    ```

    
3.  **Webhook:**

      Webhook és una forma en què una aplicació pot proporcionar informació en temps real a una altra aplicació amb poca o gens de demora. També conegut com "Reverse API".
   
      Per fer això es poden fer ús de diferents mètodes. És important que el domini o subdomini que utilitzem per webhook sigui SSL perquè si no l'API de Telegram no accepta ni peticions ni respostes a l'URL.
   
      En el nostre cas hem fet servir Ngrok.
   
      Aquest és un servei amb opcions de pagament i gratuïtes. Per poder fer l'instal·lació podeu anar a la seva pàgina web i fer les passes indicades per ells mateixos.
   
      https://ngrok.com/
    

4.  **Base de dades (master/slave):**

      Per a l'ús del nostre codi s'han d'instal·lar dues bases de dades en dues instàncies diferents. Una master i l'altre slave.
   
      Realitzar la configuració de les bases de dades en una única instància, és opcional. Per a això s'hauran de fer modificacions en el codi perquè únicament es tractin les dades des de la mateixa instància.
   
      Per a la instal·lació del motor de base de dades (En el nostre cas MariaDB):
      
      ```bash
        sudo apt install mariadb-server
      ```
   
      Per a realitzar la configuració master/slave revisar documentació del motor de base de dades:
   
      https://mariadb.com/kb/en/setting-up-replication/


5. **Prometheus**
   
   Instal·lem el servei Prometheus amb la següent comanda.
   
   ```bash
     sudo apt install prometheus -y
   ```
   Una vegada instal·lat passem a la configuració general de prometheus.

   ```bash
     sudo nano /etc/prometheus/prometheus.yml
   ```  

   En aquest arxiu podrem configurar:
   - Scrape_interval = Interval d'obtenció de dades
   - Evaluation_interval = Temps de comprovació de regles
   - Rule_files: Arxiu de regles
   - Jobs = Grups de targetes
   - Targetes = {IP:Port} Instàncies que ens interessa controlar
   
   Exemple d'arxiu configuració:

   ```yaml
   global:
     scrape_interval: 5s
     evaluation_interval: 15s
     external_labels:
       monitor: 'fuji-badalona'
   
   alerting:
     alertmanagers:
     - static_configs:
       - targets: ['192.168.1.10:9093']
   
   rule_files:
     - "node-exporter-rules.yml"
   
   scrape_configs:
     - job_name: 'prometheus'
       static_configs:
         - targets: ['192.168.1.10:9090']
   ``` 
   
   Per a l'obtenció de les mètriques en cadascuna de les instàncies s'haurà d'instal·lar l'exportador pertinent el qual volem utilitzar. En el nostre cas hem utilitzat node_exporter, mysqld_exporter i el nostre propi exportador port_exporter.

   Per a la instal·lació d'aquests exportadors pots consultar els següents enllaços:
   - Node_exporter: https://github.com/prometheus/node_exporter
   - Mysqld_exporter: https://github.com/prometheus/mysqld_exporter
   - Port_exporter: https://github.com/mohidpineda/port_exporter
   
   
6.  **Clonar Repositori:**
   
    ```bash
    git clone https://github.com/Ivit06/Telegram-Bot-With-GO.git
    ```
      Una vegada clonat el repositori s'han de canviar les credencials de l'arxiu .env amb les vostres.
      
      Per a més informació poden consultar el nostre manual tècnic:
   
      https://docs.google.com/document/d/1TidJjzkEQWyAxgp-SP6ZBXoHj5kMYe0DUvQX_fl_P2U/edit?usp=drive_link


## Execució

   Per a realitzar a l'execució del bot primer hem d'iniciar els serveis comentats, per al seu correcte funcionament:

   Prometheus:
    ```bash
    systemctl status prometheus.service
    ```   

   Ngrok:
    ```bash
    ngrok http --url=adreça_ngrok Port
    ```   

   Bot:
    ```bash
    go run main/main.go (Estant en l'arrel del projecte)
    ```   
    
   Amb tot això podríem anar a l'aplicació de Telegram i interactuar amb el nostre bot.

## Llicència
<font style="vertical-align: inherit;"><font style="vertical-align: inherit;">
    Fuji © 2025 d'Ivan Vita té llicència CC BY-NC-SA 4.0. Per veure una còpia d'aquesta llicència, visiteu https://creativecommons.org/licenses/by-nc-sa/4.0/
</font></font>

