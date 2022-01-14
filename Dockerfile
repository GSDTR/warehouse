FROM golang:1.12-stretch

RUN apt-get update
RUN apt-get install -y xvfb openjdk-8-jre unzip libgconf-2-4 chromium iceweasel bzip2
RUN go get github.com/go-redis/redis
RUN go get github.com/tebeka/selenium
RUN go get github.com/tealeg/xlsx
RUN go get github.com/lib/pq
RUN go get github.com/eclipse/paho.mqtt.golang
RUN go get github.com/eclipse/paho.mqtt.golang
#RUN go get github.com/emersion/go-imap/...
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers


ADD . /home
WORKDIR /home
RUN chmod +x vendor/geckodriver-v0.23.0-linux64
RUN chmod +x vendor/chromedriver-linux64-2.42
RUN unzip vendor/chrome-linux.zip -d vendor
RUN go build webBrowser.go
ENTRYPOINT ./webBrowser
