version: "3"
services:
  serve:
    image: "server_dev:1"
    container_name: "server"
    restart: "always"
    volumes:
      - "./Search_Engine_FrontEnd:/project"
      - "/etc/localtime:/etc/localtime"
    ports:
      - "3000:3000"
