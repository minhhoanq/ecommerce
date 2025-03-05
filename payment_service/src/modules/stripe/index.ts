import { Container } from "brandi";
import { getInstanceStripe, STRIPE_INSTANCE_TOKEN } from "./stripe";

export function bindToContainer(container: Container): void {
    container.bind(STRIPE_INSTANCE_TOKEN).toInstance(getInstanceStripe).inSingletonScope();
}