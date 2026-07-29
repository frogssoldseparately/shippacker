const fs = require("fs");
const path = require("path");

function assertExists(identifier, filepath) {
    if (!fs.existsSync(filepath)) {
        throw Error(`Directory "${filepath}" from ${identifier} arg not found.`);
    }
}

function getValidatedInput() {
    const keys = ["source", "destination"];
    const tests = [assertExists, assertExists];
    const out = new Map();
    const args = process.argv.slice(2);
    if (args.length !== keys.length) {
        throw Error(`Expected commandline arguments [${keys}]. Only ${
            args.length
        } arguments provided.`);
    }
    for (let i = 0; i < keys.length; i++) {
        const key = keys[i];
        const test = tests[i];
        const value = args[i];
        test(key, value);
        out.set(key, value);
    }
    return out;
}

function findReplacements(dir) {
    const sourceFiles = fs
        .readdirSync(dir)
        .filter(elem => elem.endsWith("Audio.txt"))
        .map(elem => path.join(dir, elem));
    const out = new Map();
    for (const filepath of sourceFiles) {
        const parts = path.basename(filepath).split("_");
        const version = parts.slice(0, 2).join("_");
        const platform = parts.slice(2, parts.length - 1).join("_");
        const key = `${version}_${platform}`;
        if (out.has(key)) {
            console.log(`Skipping duplicate replacement mapping ${filepath}`);
            continue;
        }
        const raw = fs.readFileSync(filepath, "utf8");
        out.set(key, raw.split("\n").map(elem => elem.trim()));
    }
    return out;
}

function getSamplesSection(filepath) {
    const raw = fs.readFileSync(filepath, "utf8");
    return raw.match(/<Samples.*<\/Samples>/gs);
}

function main() {
    const args = getValidatedInput();
    const srcPath = args.get("source");
    const destPath = args.get("destination");
    // Simplify and modify audio xmls
    const sourceFiles = fs
        .readdirSync(srcPath)
        .filter(elem => elem.endsWith("Audio.xml"))
        .map(elem => path.join(srcPath, elem));
    const writtenFileNames = [];
    const replacementMap = findReplacements(args.get("source"));
    for (const filepath of sourceFiles) {
        const parts = path.basename(filepath).split("_");
        const version = parts.slice(0, 2).join("_");
        const platform = parts.slice(2, 4).join("_");
        const replacements = replacementMap.get(`${version}_All`)
            ?? replacementMap.get(`${version}_${platform}`);

        const section = getSamplesSection(filepath);
        if (section == null) {
            console.log(`Couldn't parse ${filepath}. Skipping`);
            continue;
        }
        let sampleIndex = 0;
        const lines = section[0]
            .split('\n')
            .map(line => {
                line = line.trim();
                if (line.search("Sample ") === 1) {
                    const offsetAttrPosition = line.search(' Offset=');
                    if (offsetAttrPosition === -1) return '\t' + line;
                    const firstQuote = line.indexOf('"', offsetAttrPosition);
                    if (firstQuote === -1) return '\t' + line;
                    const secondQuote = line.indexOf('"', firstQuote + 1);
                    if (secondQuote === -1) return '\t' + line;
                    const offset = line.slice(firstQuote + 1, secondQuote);
                    name = replacements[sampleIndex++];
                    line = `\t<Sample Name="${name}" Offset="${offset}"/>`;
                } else if (line.search("Blob ") === 1) {
                    line = '\t' + line;
                }
                return line;
            });
        const content = `<Root>\n\t${lines.join('\n\t')}\n</Root>`;
        const outname = `${
            version
        }_${
            platform
        }_Audio_Stub.xml`;
        const outpath = path.join(destPath, outname);
        try {
            fs.writeFileSync(outpath, content, "utf8");
            writtenFileNames.push(outname);
        } catch {
            console.log(`Couldn't write to "${outpath}". Skipping`);
        }
    }
    // Update pkg/globals/audio.go
    const injections = new Map();
    const injectedKeyNames = writtenFileNames.map(elem => 
        elem.slice(0, -15).toLowerCase()
    );
    const injectedVarNames = injectedKeyNames.map(elem => `${elem}_audio_file`);
    injections.set("xml_embeds", injectedVarNames.map((elem, i) => 
        `//go:embed ${writtenFileNames[i]}\nvar ${elem} []byte`
    ).join('\n'));
    injections.set("key_var_pairs", injectedVarNames.map((elem, i) =>
        `"${injectedKeyNames[i]}" : ${elem},`
    ).join('\n\t'));
    let audioGoContent = fs.readFileSync(
        path.join(srcPath, "audio.template"),
        "utf8"
    );
    for (const [key, value] of injections) {
        audioGoContent = audioGoContent.replaceAll(`##:${key}:##`, value);
    }
    fs.writeFileSync(path.join(destPath, "audio.go"), audioGoContent, "utf8");
}

if (require.main === module) {
    main();
}