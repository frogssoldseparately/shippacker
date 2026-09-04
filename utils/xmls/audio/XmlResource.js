import fs from "fs";

export class XmlResource {
    constructor(srcPath, outPath) {
        const raw = fs.readFileSync(srcPath, "utf8");
        const section = raw.match(/<Samples.*<\/Samples>/gs);
        if (section == null) {
            throw Error("Could not find Samples.");
        }
        this.xmlEntries = section[0].split('\n').map(line => line.trim());
        this.outPath = outPath;
    }

    #prefix = "";
    set prefix(v) {
        this.#prefix = v;
    }

    #linePrefix = "";
    set linePrefix(v) {
        this.#linePrefix = v;
    }

    #tagRules = [];
    addTagRule(tagGenerator) {
        this.#tagRules.push(tagGenerator);
    }

    #suffix = "";
    set suffix(v) {
        this.#suffix = v;
    }

    writeStub() {
        const w = fs.createWriteStream(this.outPath, "utf8");
        w.write(this.#prefix);
        for (const line of this.xmlEntries) {
            w.write(this.#linePrefix);
            const tagName = line.slice(1, line.indexOf(' '));
            for (const tagRule of this.#tagRules) {
                tagRule(w, tagName, line);
            }
        }
        w.write(this.#suffix);
        w.end();
    }
}

export function extractAttribute(str, attrName) {
    const attrPosition = str.search(` ${attrName}=`);
    if (attrPosition === -1) return null;
    const firstQuote = str.indexOf('"', attrPosition);
    const secondQuote = str.indexOf('"', firstQuote + 1);
    return str.slice(firstQuote + 1, secondQuote);
}