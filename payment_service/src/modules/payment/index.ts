import { Container } from "brandi";
import { PAYMENT_MANAGEMENT_OPERATOR_IMPL_TOKEN, PaymentManagementOperatorImpl } from "./payment_management_operators";

export function bindToContainer(container: Container): void {
    container.bind(PAYMENT_MANAGEMENT_OPERATOR_IMPL_TOKEN).toInstance(PaymentManagementOperatorImpl).inSingletonScope();
}