import { Container } from "brandi";
import { LOGGER_WINSTON_TOKEN, LoggerWinston } from "./logger";
import { BINARY_CONVERTOR_TOKEN, BinaryConvertorImpl } from "./binary_converter";

export * from "./logger";
export * from "./errors";
export * from "./binary_converter"

export function bindToContainer(container: Container): void {
    container.bind(LOGGER_WINSTON_TOKEN).toInstance(LoggerWinston).inSingletonScope();
    container.bind(BINARY_CONVERTOR_TOKEN).toInstance(BinaryConvertorImpl).inSingletonScope();
}
