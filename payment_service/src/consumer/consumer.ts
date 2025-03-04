import { injected, token } from "brandi";
import { MESSAGE_CONSUMER_TOKEN, MessageConsumer } from "../dataaccess/kafka/consumer";
import { CreatePaymentRequest } from "../proto/gen/payment_service/CreatePaymentRequest";
import { BINARY_CONVERTOR_TOKEN, BinaryConvertor, LOGGER_WINSTON_TOKEN, LoggerWinston } from "../utils";
import { PAYMENT_CREATED_MESSAGE_HANDLER_IMPL_TOKEN, PaymentCreatedMessageHandler } from "./payment_created";

const TOPIC_NAME_ORDER_SERVICE_ORDER_CREATED = "order_serive_order_created"

export class PaymentKafkaConsumer {
    constructor(
        private readonly messageConsumer: MessageConsumer,
        private readonly paymentCreatedMessageHandler: PaymentCreatedMessageHandler,
        private readonly binaryConvertor: BinaryConvertor,
        private readonly logger: LoggerWinston
    ) { }

    public start(): void {
        this.messageConsumer.registerHandlerEvent([{
            topic: TOPIC_NAME_ORDER_SERVICE_ORDER_CREATED,
            onMessage: (message) =>
                this.onPaymentCreated(message),
        }])
            .then(() => {
                if (process.send) {
                    process.send("ready")
                }
            })
    }

    private async onPaymentCreated(message: Buffer | null): Promise<void> {
        if (message === null) {
            this.logger.error("null message, skipping")
            return;
        }

        const paymentCreatedMessage = this.binaryConvertor.fromBuffer(message);
        // TODO
        await this.paymentCreatedMessageHandler.paymentCreated(paymentCreatedMessage);
    }
}

injected(PaymentKafkaConsumer, MESSAGE_CONSUMER_TOKEN, PAYMENT_CREATED_MESSAGE_HANDLER_IMPL_TOKEN, BINARY_CONVERTOR_TOKEN, LOGGER_WINSTON_TOKEN)

export const PAYMENT_KAFKA_CONSUMER_TOKEN = token<PaymentKafkaConsumer>("PaymentKafkaConsumer");
