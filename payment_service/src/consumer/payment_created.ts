import { injected, token } from "brandi";
import { PAYMENT_MANAGEMENT_OPERATOR_IMPL_TOKEN, PaymentManagementOperator } from "../modules/payment/payment_management_operators";
import { CreatePaymentRequest } from "../proto/gen/payment_service/CreatePaymentRequest";
import { LOGGER_WINSTON_TOKEN, LoggerWinston } from "../utils";

export interface PaymentCreatedMessageHandler {
    paymentCreated(message: CreatePaymentRequest): Promise<void>
}

export class PaymentCreatedMessageHandlerImpl implements PaymentCreatedMessageHandler {
    constructor(
        private readonly paymentManagementOperator: PaymentManagementOperator,
        private readonly logger: LoggerWinston
    ) { }

    public async paymentCreated(message: CreatePaymentRequest): Promise<void> {
        this.logger.info(`payment_service_payment_created message received: ${message}`)

        if (message.orderId === "") {
            this.logger.error(`payment_service_payment_created order_id is required: ${message.orderId}`)
            return;
        }

        await this.paymentManagementOperator.createPayment(message)
    }
}

injected(PaymentCreatedMessageHandlerImpl, PAYMENT_MANAGEMENT_OPERATOR_IMPL_TOKEN, LOGGER_WINSTON_TOKEN)

export const PAYMENT_CREATED_MESSAGE_HANDLER_IMPL_TOKEN = token<PaymentCreatedMessageHandlerImpl>("PaymentCreatedMessageHandlerImpl")
