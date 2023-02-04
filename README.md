

## QoS LEVEL 2

  - kuchat init
      - return uuid as seesion_id

  
  - kuchat <session_id> <user_email> <destination_email>
      - start ui
      - subscribe to topic `/chats/<session_id>/<user_email>`
      - send to topic `/chats/<session_id>/<destination_email>`


