import { Producer } from "kafkajs";
import { PaymentTransactionStatus } from "../../db/models";
import { TopicOrderPaymentCreated } from "../../../utils/constants";
import { ErrorWithStatus, LOGGER_WINSTON_TOKEN, LoggerWinston } from "../../../utils";
import { status } from "@grpc/grpc-js";
import { KAFKA_PRODUCER_TOKEN } from "./producer";
import { injected, token } from "brandi";

export class PaymentTransactionCompleted {
    constructor(
        public orderId: string,
        public paymentTransactionStatus: PaymentTransactionStatus
    ) { }
}

export interface PaymentTransactionCompletedProducer {
    createPaymentTransactionCompletedMessage(message: PaymentTransactionCompleted): Promise<void>
}

export class PaymentTransactionCompletedProducerImpl implements PaymentTransactionCompletedProducer {
    constructor(
        private readonly producer: Producer,
        private readonly logger: LoggerWinston
    ) { }

    public async createPaymentTransactionCompletedMessage(message: PaymentTransactionCompleted): Promise<void> {
        try {
            await this.producer.connect()
            await this.producer.send({
                topic: TopicOrderPaymentCreated,
                messages: [{ value: "" }]
            }).then(() => this.logger.info("send message payment transaction success"))
        } catch (error) {
            this.logger.error(
                `failed to create ${TopicOrderPaymentCreated} message with error: ${error}`)
            throw ErrorWithStatus.withStatus(error, status.INTERNAL)
        }
    }
}

injected(PaymentTransactionCompletedProducerImpl, KAFKA_PRODUCER_TOKEN, LOGGER_WINSTON_TOKEN);

export const PAYMENT_TRANSACTION_COMPLETED_PRODUCER_TOKEN = token<PaymentTransactionCompletedProducerImpl>("PaymentTransactionCompletedProducerImpl");
