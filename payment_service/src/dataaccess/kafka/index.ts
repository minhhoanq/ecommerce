import { Container } from "brandi";
import { getInstanceKafka, KAFKA_INSTANCE_TOKEN } from "./kafka";
import * as consumer from "./consumer"
import * as producer from "./producer"

export function bindToContainer(container: Container): void {
    container.bind(KAFKA_INSTANCE_TOKEN).toInstance(getInstanceKafka).inSingletonScope();
    producer.bindToContainer(container);
    consumer.bindToContainer(container);
}
