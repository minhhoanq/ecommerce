import { Container } from "brandi";
import { getKafkaProducer, KAFKA_PRODUCER_TOKEN } from "./producer";
import { PAYMENT_TRANSACTION_COMPLETED_PRODUCER_TOKEN, PaymentTransactionCompletedProducerImpl } from "./payment_transaction_completed";

export * from "./producer";
export * from "./payment_transaction_completed"

export function bindToContainer(container: Container): void {
    container.bind(KAFKA_PRODUCER_TOKEN).toInstance(getKafkaProducer).inSingletonScope();
    container.bind(PAYMENT_TRANSACTION_COMPLETED_PRODUCER_TOKEN).toInstance(PaymentTransactionCompletedProducerImpl).inSingletonScope();
}
