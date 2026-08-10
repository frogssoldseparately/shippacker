import { join } from "path";
import { fileURLToPath } from "url";
import { getValidatedInput } from "./audio/validate.js";
import { makeStubs, fillAudioTemplate, fillSetupTemplate } from "./audio/stubs.js";
import { parseSetups } from "./audio/SetupEntry.js";

/**
 * Generates abridged Audio.xml files with names matching mm.o2r and embeds them
 * into Ship Packer.
 */
function main() {
    const cmdArgs = getValidatedInput();
    const source = cmdArgs.get("source");
    const destination = cmdArgs.get("destination");

    const audioTemplatePath = join(source, "audio.template");
    const audioGoPath = join(destination, "audio.go");
    const setupTemplatePath = join(source, "setup.template");
    const setupJsonPath = join(source, "setup.json");
    const setupGoPath = join(destination, "setup.go");

    const versionSetups = parseSetups(setupJsonPath);

    fillAudioTemplate(
        audioTemplatePath, audioGoPath, makeStubs(source, destination), versionSetups
    );
    fillSetupTemplate(
        setupTemplatePath, setupGoPath, versionSetups
    )
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
    main();
}