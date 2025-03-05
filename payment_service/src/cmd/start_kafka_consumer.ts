import { Container } from "brandi"
import * as config from "../config"
import * as db from "../dataaccess/db"
import * as kafka from "../dataaccess/kafka"
import * as utils from "../utils"
import * as service from "../service"
import * as modules from "../modules"
import * as consumer from "../consumer"
import dotenv from "dotenv"

export async function startKafkaConsumer(dotenvPath: string): Promise<void> {
    dotenv.config({
        path: dotenvPath
    })

    const container = new Container();
    kafka.bindToContainer(container);
    config.bindToContainer(container);
    db.bindToContainer(container);
    utils.bindToContainer(container);
    modules.bindToContainer(container);
    service.bindToContainer(container);
    consumer.bindToContainer(container);

    const kafkaConsumer = container.get(consumer.PAYMENT_KAFKA_CONSUMER_TOKEN);
    kafkaConsumer.start();
}
