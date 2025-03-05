import Stripe from "stripe";
import { STRIPE_CONFIG_TOKEN, StripeConfig } from "../../config/stripe";
import { injected, token } from "brandi";

export function getInstanceStripe(config: StripeConfig): Stripe {
    return new Stripe(config.apiKey)
}

injected(getInstanceStripe, STRIPE_CONFIG_TOKEN)

export const STRIPE_INSTANCE_TOKEN = token<Stripe>("Stripe");
