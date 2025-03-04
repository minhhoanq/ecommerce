import { Payment as PaymentProto } from "../proto/gen/payment_service/Payment";

export enum PaymentStatus {
    PENDING = "PENDING",
    SUCCESS = "SUCCESS",
    FAILED = "FAILED",
}

export enum PaymentMethod {
    CASH = "cash",
    CREDIT_CARD = "credit card",
}

export class Payment {
    constructor(
        public orderId: string,
        public onl_payment_intent_id: string | null,
        public amount: number,
        public status: PaymentStatus,
        public paymentMethod: PaymentMethod
    ) {
        this.status = PaymentStatus.PENDING;
        this.paymentMethod = PaymentMethod.CASH;
    }

    public fromProto(paymentProto: PaymentProto | undefined | null): Payment {
        return new Payment(
            paymentProto?.id || "",
            paymentProto?.onlPaymentIntentId || null,
            paymentProto?.amount as number || 0,
            paymentProto?.status as PaymentStatus || PaymentStatus.PENDING,
            paymentProto?.paymentMethod as PaymentMethod || PaymentMethod.CASH
        )
    }
}