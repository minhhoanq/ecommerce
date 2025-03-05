import { Container } from "brandi"
import dotenv from "dotenv"
import * as config from "../config";
import * as utils from "../utils";
import * as service from "../service";
import * as payment from "../modules/payment";
import * as db from "../dataaccess/db";

export async function startGRPCServer(dotenvPath: string): Promise<void> {
    dotenv.config({
        path: dotenvPath
    })

    const container = new Container();
    utils.bindToContainer(container);
    config.bindToContainer(container);
    service.bindToContainer(container);
    payment.bindToContainer(container);
    db.bindToContainer(container);

    const server = container.get(service.PAYMENT_SERVICE_GRPC_SERVER_TOKEN);
    server.loadProtoAndStartServer("./src/proto/payment_service.proto")
}
