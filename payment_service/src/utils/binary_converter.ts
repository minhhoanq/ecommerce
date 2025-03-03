import { injected, token } from "brandi";

export interface BinaryConvertor {
    toBuffer<T = any>(data: T): Buffer;
    fromBuffer<T = any>(buffer: Buffer): T | null;
}

export class BinaryConvertorImpl implements BinaryConvertor {
    toBuffer<T = any>(data: T): Buffer {
        return Buffer.from(JSON.stringify(data));
    }

    fromBuffer<T = any>(buffer: Buffer): T | null {
        const jsonStr = buffer.toString()
        if (jsonStr === "") {
            return null
        }
        return JSON.parse(jsonStr);
    }
}

injected(BinaryConvertorImpl);

export const BINARY_CONVERTOR_TOKEN = token<BinaryConvertorImpl>("BinaryConvertorImpl");
