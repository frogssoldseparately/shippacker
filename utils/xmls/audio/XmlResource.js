import fs from "fs";

export class XmlResource {
    constructor(srcPath, replacements, outPath) {
        const raw = fs.readFileSync(srcPath, "utf8");
        const section = raw.match(/<Samples.*<\/Samples>/gs);
        if (section == null) {
            throw Error("Could not find Samples.");
        }
        this.xmlEntries = section[0].split('\n').map(line => line.trim());
        this.replacements = replacements;
        this.outPath = outPath;
    }

    writeStub() {
        const w = fs.createWriteStream(this.outPath, "utf8");
        w.write('<Root>');
        let sampleIndex = 0;
        for (const line of this.xmlEntries) {
            w.write('\n');
            const firstSpace = line.indexOf(" ");
            const tagName = line.slice(1, firstSpace);
            if (tagName === "Sample") {
                w.write('\t\t');
                const offsetAttrPosition = line.search(' Offset=');
                const firstQuote = line.indexOf('"', offsetAttrPosition);
                const secondQuote = line.indexOf('"', firstQuote + 1);
                const offset = line.slice(firstQuote + 1, secondQuote);
                if (offset.length === 0) {
                    w.write(line);
                } else {
                    const name = this.replacements[sampleIndex++];
                    w.write(`<Sample Name="${name}" Offset="${offset}"/>`);
                }
            } else {
                if (tagName === "Blob") {
                    w.write('\t');
                }
                w.write('\t');
                w.write(line);
            }
        }
        w.write('\n</Root>');
        w.end();
    }
}