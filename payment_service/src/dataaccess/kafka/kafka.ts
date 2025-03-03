import { injected, token } from "brandi";
import { Kafka, KafkaConfig } from "kafkajs";
import { KAFKA_CONFIG_TOKEN } from "../../config/kafka";

export function getInstanceKafka(config: KafkaConfig): Kafka {
    return new Kafka({
        clientId: config.clientId,
        brokers: config.brokers
    })
}

injected(getInstanceKafka, KAFKA_CONFIG_TOKEN)

export const KAFKA_INSTANCE_TOKEN = token<Kafka>("Kafka")
