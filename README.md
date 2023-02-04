

## Chat Service :
  - Register user -> generate empty acl
      - endpoint POST `/register`
      - body:
          - email
          - password

  - Login get credential
      - endpoint POST `/login`
      - body:
          - email
          - password
      - response credential format `<email>:<client_id>:<acl_password>`::base64

  - get sessions
      - endpoint GET `/sessions`
      - header basic auth email:password

  - new sessions
      - endpoint POST `/sessions`
      - header basic auth email:password

  - new session
      - endpoint POST `/session/:receiver_email`
      - header basic auth email:password
      - will check session exist, create if not exist
      - if session exist go to start session
      - if session doesnt exist, create session -> add/update both users acl's
      - return :
          - subscribe topic
          - publish topic
          - info topic

  - start session
      - endpoint POST `/session/:session_id/start`
      - header basic auth email:password
      - return :
          - subscribe topic
          - publish topic
          - info topic



## QoS LEVEL 1

## Chat client

  - kuchat user create <new_email> <password>

  
  - kuchat user login <email> <password>
      - response : credential save to .kuchat


  - kuchat sessions list
      - list all sessions

  - kuchat session init <receiver email>
      - generate new chat session start ui

  - kuchat session start <session_id>
      - get topics and start pub sub open chat ui



## NEXT TODO CHAT HISTORY



