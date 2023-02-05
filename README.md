## Simple MQTT Based Realtime Chat TUI

<img width="638" alt="Screenshot 2023-02-05 at 10 05 09" src="https://user-images.githubusercontent.com/40946917/216800041-41c99693-5b7f-41b1-9e24-680f0d0ede5f.png">

## Article :
  - i wrote article about this https://alfiankan.medium.com/realtime-chat-tanpa-web-socket-dengan-mqtt-protocol-3108898e51b0


## Need To Read :
  - docker compose using arm image by default, replace using general vernemq image tag.
  
  ```yaml
  version: "3.7"

services:
  vmq:
    image: vernemq/vernemq
    ports:
      - 8080:8080
      - 8888:8888
      - 1883:1883
    volumes:
      - ./etc:/vernemq/etc/
    environment:
      DOCKER_VERNEMQ_ACCEPT_EULA: yes
  ```

  - Dashboard status `http://localhost:8888/status`
  - QoS level Two

## How To Use

  - kuchat init
      - return uuid as seesion_id

  
  - kuchat <session_id> <user_email> <destination_email>
      - start ui
      - subscribe to topic `/chats/<session_id>/<user_email>`
      - send to topic `/chats/<session_id>/<destination_email>`


