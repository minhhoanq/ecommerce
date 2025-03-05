export class Payment {
    constructor(
        public id: string,
        public order_id: string,
        public onl_payment_intent_id: string,
        public amount: number,
        public status: PaymentTransactionStatus,
        public payment_method: string,
    ) { }
}

export enum PaymentTransactionStatus {
    PENDING = "PENDING",
    SUCCESS = "SUCCESS",
    CANCEL = "CANCLE"
}
