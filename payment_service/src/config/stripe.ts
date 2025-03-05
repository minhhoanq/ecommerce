import { token } from "brandi";

export class StripeConfig {
    public apiKey = "";
    public webHookEndpointSecret = "";
    public port = 4242;

    public static fromEnv(): StripeConfig {
        const config = new StripeConfig();
        config.apiKey = process.env.STRIPE_API_KEY!;
        config.webHookEndpointSecret = process.env.STRIPE_WEBHOOK_ENDPOINT_SECRET!;
        config.port = +process.env.port!;
        return config;
    }
}

export const STRIPE_CONFIG_TOKEN = token<StripeConfig>("StripeConfig");
