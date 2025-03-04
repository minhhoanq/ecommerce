import { Container } from "brandi";
import { PAYMENT_CREATED_MESSAGE_HANDLER_IMPL_TOKEN, PaymentCreatedMessageHandlerImpl } from "./payment_created";
import { PAYMENT_KAFKA_CONSUMER_TOKEN, PaymentKafkaConsumer } from "./consumer";

export function bindToContainer(container: Container): void {
    container.bind(PAYMENT_KAFKA_CONSUMER_TOKEN).toInstance(PaymentKafkaConsumer).inSingletonScope();
    container.bind(PAYMENT_CREATED_MESSAGE_HANDLER_IMPL_TOKEN).toInstance(PaymentCreatedMessageHandlerImpl).inSingletonScope();
}