import { status } from "@grpc/grpc-js";
import { Payment } from "../../dataaccess/db/models";
import { CreatePaymentParams, PAYMENT_DATA_ACCESSOR_TOKEN, PaymentDataAccessor } from "../../dataaccess/db/payment";
import { CreatePaymentRequest } from "../../proto/gen/payment_service/CreatePaymentRequest";
import { CreatePaymentResponse } from "../../proto/gen/payment_service/CreatePaymentResponse";
import { ErrorWithStatus } from "../../utils";
import { injected, token } from "brandi";
import { Payment as paymentRes } from "../../proto/gen/payment_service/Payment";

export interface PaymentManagementOperator {
    createPayment(arg: CreatePaymentRequest): Promise<CreatePaymentResponse>;
}

export class PaymentManagementOperatorImpl implements PaymentManagementOperator {
    constructor(
        private readonly paymentDataAccessor: PaymentDataAccessor
    ) { }

    public async createPayment(arg: CreatePaymentRequest): Promise<CreatePaymentResponse> {
        const params: CreatePaymentParams = {
            order_id: arg.orderId ?? "",
            onl_payment_intent_id: arg.onlPaymentIntentId ?? "",
            amount: Number(arg.amount) ?? 0,
            payment_method: arg.paymentMethod ?? ""
        }
        const payment: Payment | null = await this.paymentDataAccessor.createPayment(params);

        if (payment == null) {
            throw new ErrorWithStatus(`failed to create payment`, status.INTERNAL)
        }

        const paymentRes: paymentRes = {
            id: payment.id,
            orderId: payment.order_id,
            onlPaymentIntentId: payment.onl_payment_intent_id,
            amount: payment.amount,
            status: 'PENDING',
            paymentMethod: payment.payment_method,
        }

        const response: CreatePaymentResponse = {
            payment: paymentRes
        }

        return response
    }
}

injected(PaymentManagementOperatorImpl, PAYMENT_DATA_ACCESSOR_TOKEN);

export const PAYMENT_MANAGEMENT_OPERATOR_IMPL_TOKEN = token<PaymentManagementOperatorImpl>("PaymentManagementOperatorImpl");
