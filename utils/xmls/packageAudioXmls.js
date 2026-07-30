import { join } from "path";
import { fileURLToPath } from "url";
import { getValidatedInput } from "./audio/validate.js";
import { makeStubs, fillTemplate } from "./audio/stubs.js";

/**
 * Generates abridged Audio.xml files with names matching mm.o2r and embeds them
 * into Ship Packer.
 */
function main() {
    const cmdArgs = getValidatedInput();
    const source = cmdArgs.get("source");
    const destination = cmdArgs.get("destination");
    fillTemplate(
        join(source, "audio.template"),
        join(destination, "audio.go"),
        makeStubs(source, destination)
    );
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
    main();
}