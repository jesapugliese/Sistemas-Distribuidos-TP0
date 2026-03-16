# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

### Ejercicio N°1:

Este ejercicio consiste en la definición de un script de bash `generar-compose.sh` que permite crear una definición de Docker Compose con una cantidad configurable de clientes.

#### Ejecución:

1. Darle permisos de ejecución al archivo:  
   ```bash
   chmod +x generar-compose.sh
   ```

2. Ejecutar el programa:  
   ```bash
   ./generar-compose.sh <nombre_archivo> <numero_clientes>
   ```

   - `<nombre_archivo>`: Nombre del archivo de salida.
   - `<numero_clientes>`: Número de clientes.

   Ejemplo: 
   ```bash
   ./generar-compose.sh docker-compose-dev.yaml 5
   ```

### Ejercicio N°2:

Para lograr que realizar cambios en los archivos de configuración del cliente y el servidor no requiera reconstruir las imágenes de Docker para que los mismos sean efectivos, se utilizaron volumenes en cada container de tal forma que el archivo de configuración usado sea el mismo al que tiene acceso el Host OS.  
Por ejemplo, para el caso del server:

```yaml
volumes:
  - ./server/config.ini:/config.ini
```

De esta forma tanto el filesystem del container y el Host OS trabajan sobre el mismo archivo.  
A su vez se evitó copiar el archivo de configuración en la imagen. Esto se logró editando el comando `COPY` dentro de los archivos Dockerfile del cliente y el servidor. Esta última modificación era innecesaria (pues un volumen tiene prioridad sobre el filesystem de la imagen si se monta en el mismo path), pero es más limpio y eficiente al copiar únicamente los archivos necesarios.  
Finalmente, se eliminó la configuración de las variables de entorno de _log levels_ del YAML para permitir que se tomen aquellas definidas en los archivos de configuración. 
