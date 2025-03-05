import { Knex } from "knex"
import { ErrorWithStatus, LOGGER_WINSTON_TOKEN, LoggerWinston } from "../../utils"
import { Payment } from "./models"
import { status } from "@grpc/grpc-js"
import {
    COL_PAYMENT_AMOUNT,
    COL_PAYMENT_ID,
    COL_PAYMENT_ONL_PAYMENT_INTENT_ID,
    COL_PAYMENT_ORDER_ID,
    COL_PAYMENT_PAYMENT_METHOD,
    COL_PAYMENT_STATUS
} from "./constants"
import { injected, token } from "brandi"
import { KNEX_INSTANCE_TOKEN } from "./knex"

export interface CreatePaymentParams {
    order_id: string,
    onl_payment_intent_id: string | null,
    amount: number,
    payment_method: string
}

export interface PaymentDataAccessor {
    createPayment(arg: CreatePaymentParams): Promise<Payment | null>
}

export class PaymentDataAccessorImpl implements PaymentDataAccessorImpl {
    constructor(
        private readonly logger: LoggerWinston,
        private readonly knex: Knex<any, any[]>
    ) {
        this.logger = new LoggerWinston()
    }

    public async createPayment(arg: CreatePaymentParams): Promise<Payment | null> {
        try {
            this.logger.info(`starting to create payment: ${arg.amount}`)
            this.logger.info(`starting to create payment: ${arg.order_id}`)
            this.logger.info(`starting to create payment: ${arg.onl_payment_intent_id}`)
            this.logger.info(`starting to create payment: ${arg.payment_method}`)

            const rows = await this.knex
                .insert({
                    [COL_PAYMENT_ORDER_ID]: arg.order_id,
                    [COL_PAYMENT_ONL_PAYMENT_INTENT_ID]: arg.onl_payment_intent_id,
                    [COL_PAYMENT_AMOUNT]: arg.amount,
                    [COL_PAYMENT_PAYMENT_METHOD]: arg.payment_method,
                })
                .returning("*")
                .into("payments");

            const response: Payment = new Payment(
                rows[0][COL_PAYMENT_ID],
                rows[0][COL_PAYMENT_ORDER_ID],
                rows[0][COL_PAYMENT_ONL_PAYMENT_INTENT_ID],
                rows[0][COL_PAYMENT_AMOUNT],
                rows[0][COL_PAYMENT_STATUS],
                rows[0][COL_PAYMENT_PAYMENT_METHOD]
            )

            return response;
        } catch (error) {
            this.logger.error(`failed to create payment: ${error}`)
            throw ErrorWithStatus.withStatus(error, status.INTERNAL)
        }
    }
}

injected(PaymentDataAccessorImpl, LOGGER_WINSTON_TOKEN, KNEX_INSTANCE_TOKEN)

export const PAYMENT_DATA_ACCESSOR_TOKEN = token<PaymentDataAccessorImpl>("PaymentDataAccessorImpl")
