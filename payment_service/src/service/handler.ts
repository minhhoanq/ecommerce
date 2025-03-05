import { sendUnaryData, status } from "@grpc/grpc-js";
import { PaymentServiceHandlers } from "../proto/gen/payment_service/PaymentService";
import { ErrorWithStatus } from "../utils";
import { CreatePaymentResponse } from "../proto/gen/payment_service/CreatePaymentResponse";
import { injected, token } from "brandi";
import { PAYMENT_MANAGEMENT_OPERATOR_IMPL_TOKEN, PaymentManagementOperator } from "../modules/payment/payment_management_operators";

export class PaymentServiceHandlerFactory {
    constructor(
        private readonly paymentManagementOperators: PaymentManagementOperator
    ) { }

    public handlers(): PaymentServiceHandlers {
        const handler: PaymentServiceHandlers = {
            CreatePayment: async (call, callback) => {
                const req = call.request;
                const response: CreatePaymentResponse = await this.paymentManagementOperators.createPayment(req)

                callback(null, response)
            }
        }

        return handler
    }

    public handlerError(error: unknown, callback: sendUnaryData<any>) {
        if (error instanceof ErrorWithStatus) {
            callback({ message: error.message, code: error.status });
        } else if (error instanceof Error) {
            callback({ message: error.message, code: status.INTERNAL });
        } else {
            callback({ code: status.INTERNAL });
        }
    }
}

injected(PaymentServiceHandlerFactory, PAYMENT_MANAGEMENT_OPERATOR_IMPL_TOKEN)

export const PAYMENT_SERVICE_HANDLER_FACTORY_TOKEN = token<PaymentServiceHandlerFactory>("PaymentServiceHandlerFactory");
