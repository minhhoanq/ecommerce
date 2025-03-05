import { Container } from "brandi";
import * as payment from "./payment";

export function bindToContainer(container: Container): void {
    payment.bindToContainer(container);
}
