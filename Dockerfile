FROM golang:1.18.1-alpine AS builder

WORKDIR /src

COPY /envoyer_backend/go.mod .
COPY /envoyer_backend/go.sum .

RUN go mod download

COPY /envoyer_backend .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app .

FROM ubuntu:latest AS final

# Install wget
RUN apt update && apt install -y build-essential wget

#install rabbitmq
RUN wget http://archive.ubuntu.com/ubuntu/pool/main/r/rabbitmq-server/rabbitmq-server_3.10.8-1.1_all.deb
RUN apt-get install -y ./rabbitmq-server_3.10.8-1.1_all.deb --fix-missing
RUN apt-get update
RUN rm ./rabbitmq-server_3.10.8-1.1_all.deb


# install curl, nodejs and yarn
RUN apt-get update && \
    apt-get install -y curl && \
    curl -sL https://deb.nodesource.com/setup_16.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g yarn

#install mysql database
RUN apt-get install -y mysql-server

# save rabbitmq and database data
VOLUME /var/lib/mysql
VOLUME /var/lib/rabbitmq

#run mysql server, create user and give permission to access outside the localhost. user = newuser, password = root
RUN service mysql start && sleep 5 && \
    echo "CREATE DATABASE IF NOT EXISTS envoyer;" | mysql -u root --password=root &&  \
    echo "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';" | mysql -u root --password=root &&  \
    echo "CREATE USER 'newuser'@'%' IDENTIFIED BY 'password';" | mysql -u root --password=root &&  \
    echo "GRANT ALL PRIVILEGES ON envoyer.* TO 'newuser'@'%';" | mysql -u root --password=root  &&  \
    echo "ALTER USER 'newuser'@'%' IDENTIFIED WITH mysql_native_password BY 'root';" | mysql -u root --password=root && \
    echo "FLUSH PRIVILEGES;" | mysql -u root --password=root && \
    service mysql stop
RUN echo "[mysqld]" >> /etc/mysql/my.cnf
RUN echo "bind-address=0.0.0.0" >> /etc/mysql/my.cnf

WORKDIR /src

# Import the compiled executable from the first stage.
COPY --from=builder /src/base.env .
COPY --from=builder /src/.env .
COPY --from=builder /app .
COPY --from=builder /src/localize ./localize
COPY /run.sh .

#copy plugin
COPY /lib/rabbitmq_delayed_message_exchange-3.10.2.ez /usr/lib/rabbitmq/lib/rabbitmq_server-3.10.8/plugins/

#enable plugins
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange
RUN rabbitmq-plugins enable rabbitmq_management

#setup frontend
WORKDIR /src/envoyer_frontend
COPY /envoyer_frontend .
RUN yarn install --frozen-lockfile
RUN yarn build
WORKDIR /src

#backend 8081, rabbitmq management 15672, mysql database 3306, frontend 3000
EXPOSE 8081 15672 3306 3000

RUN echo "NODENAME=rabbit@localhost" > /etc/rabbitmq/rabbitmq-env.conf
ENV RABBITMQ_SERVER_START_ARGS -eval error_logger:tty(true).

# add new user admin admin for rabbitmq
RUN rabbitmq-server -detached && sleep 10 &&  \
    rabbitmqctl add_user admin admin &&  \
    rabbitmqctl set_permissions -p / admin ".*" ".*" ".*" &&  \
    rabbitmqctl set_user_tags admin administrator &&  \
    rabbitmqctl stop

#give execute permission on run.sh
RUN chmod +x ./run.sh

CMD ["./run.sh"]


