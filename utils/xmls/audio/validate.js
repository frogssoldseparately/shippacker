import fs from "fs";

function assertExists(identifier, filepath) {
    if (!fs.existsSync(filepath)) {
        throw Error(`Directory "${filepath}" from ${identifier} arg not found.`);
    }
}

export function getValidatedInput() {
    const keys = ["source", "destination"];
    const tests = [assertExists, assertExists];
    const out = new Map();
    const cmdArgs = process.argv.slice(2);
    if (cmdArgs.length !== keys.length) {
        throw Error(`Expected commandline arguments [${keys}]. Only ${
            cmdArgs.length
        } arguments provided.`);
    }
    for (let i = 0; i < keys.length; i++) {
        const key = keys[i];
        const test = tests[i];
        const value = cmdArgs[i];
        test(key, value);
        out.set(key, value);
    }
    return out;
}