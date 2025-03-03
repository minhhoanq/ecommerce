import { Consumer } from "kafkajs";
import { KAFKA_CONFIG_TOKEN, KafkaConfig } from "../../../config/kafka";
import { ErrorWithStatus, LOGGER_WINSTON_TOKEN, LoggerWinston } from "../../../utils";
import { injected, token } from "brandi";
import { KAFKA_CONSUMER_TOKEN } from "./consumer";

export type MessageHandlerFunc = (message: Buffer | null) => Promise<void>;

export interface MessageHandler {
    topic: string;
    onMessage: MessageHandlerFunc;
}

export interface MessageConsumer {
    registerHandlerEvent(handlerList: MessageHandler[]): Promise<void>;
}

export class MessageConsumerImpl implements MessageConsumer {
    constructor(
        private readonly consumer: Consumer,
        private readonly kafkaConfig: KafkaConfig,
        private readonly logger: LoggerWinston
    ) { }

    public async registerHandlerEvent(handlerList: MessageHandler[]): Promise<void> {
        const topicToHandlerMap = this.getTopicToHandlerMap(handlerList)
        await this.consumer.connect()
        for (const handler of handlerList) {
            await this.consumer.subscribe({
                topic: handler.topic,
                fromBeginning: true
            })
        }

        await this.consumer.run({
            eachMessage: async ({ topic, message, heartbeat }) => {
                const handler = topicToHandlerMap.get(topic);
                if (handler === undefined) {
                    return;
                }

                const heartbeatInterval = setInterval(
                    async () => await heartbeat(),
                    this.kafkaConfig.heartbeatInterval
                )

                try {
                    handler(message.value)
                } catch (error) {
                    this.logger.error(`failed to handle message topic: ${topic}, message: ${message.value}, error: ${error}`)
                    throw error;
                } finally {
                    clearInterval(heartbeatInterval)
                }
            }
        })
    }

    private getTopicToHandlerMap(handlerList: MessageHandler[]): Map<string, MessageHandlerFunc> {
        return new Map(handlerList.map(handler => [handler.topic, handler.onMessage]))
    }
}

injected(MessageConsumerImpl, KAFKA_CONSUMER_TOKEN, KAFKA_CONFIG_TOKEN, LOGGER_WINSTON_TOKEN)

export const MESSAGE_CONSUMER_TOKEN = token<MessageConsumerImpl>("MessageConsumerImpl")
