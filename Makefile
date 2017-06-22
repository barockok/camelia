LIBRDKAFKA_PATH = /tmp/.lib-rdkafka-build
clean_librdkafka:
	rm -rf ${LIBRDKAFKA_PATH}
	rm -rf ${GOPATH}/pkg/darwin_amd64/github.com/confluentinc/confluent-kafka-go
	rm -rf ${GOPATH}/src/github.com/confluentinc/confluent-kafka-go
build_librdkafka: clean_librdkafka
	mkdir -p ${LIBRDKAFKA_PATH}
	cd ${LIBRDKAFKA_PATH} ; \
	test ! -d ${LIBRDKAFKA_PATH}/librdkafka && git clone https://github.com/edenhill/librdkafka.git
	cd ${LIBRDKAFKA_PATH}/librdkafka ; \
	./configure && make && make install
	go get -u github.com/confluentinc/confluent-kafka-go/kafka