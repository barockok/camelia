LIBRDKAFKA_PATH = /tmp/onta-kafka-build
build_librdkafka:
	mkdir -p ${LIBRDKAFKA_PATH}
	cd ${LIBRDKAFKA_PATH} ; \
	test ! -d ${LIBRDKAFKA_PATH}/librdkafka && git clone https://github.com/edenhill/librdkafka.git
	cd ${LIBRDKAFKA_PATH}/librdkafka ; \
	./configure --prefix ${LIBRDKAFKA_PATH}/usr && make && make install
	PKG_CONFIG_PATH=${LIBRDKAFKA_PATH}/usr/lib/pkgconfig go get -u github.com/confluentinc/confluent-kafka-go/kafka 
	rm -rf ${LIBRDKAFKA_PATH}