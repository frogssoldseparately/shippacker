import fs from "fs";

export class SetupEntry {
    constructor(name, obj) {
        this.name = name;
        this.assignStrings = [];
        if (obj.hasOwnProperty("Settings")) {
            for (const settingName of Object.getOwnPropertyNames(obj.Settings)) {
                let val = obj.Settings[settingName];
                switch (typeof(val)) {
                    case "string":
                        val = `"${val}"`;
                        break;
                    default:
                        // do nothing
                }
                this.assignStrings.push(`${settingName} = ${val}`)
            }
        }
        if (obj.hasOwnProperty("WarnUnstable") && obj.WarnUnstable) {
            this.assignStrings.push("if err := iohelper.WarnUnstable(); " +
                "err != nil {\n\t\t\treturn err\n\t\t}");
        }
    }
}

/**
 * Parses setup config information.
 * @param {string} srcPath - setup.json path.
 * @returns Array of SetupEntry instances.
 */
export function parseSetups(srcPath) {
    if (!fs.existsSync(srcPath)) {
        console.log(`Setup:Bad: Could not read ${
            srcPath
        }. It does not exist.`);
        return [];
    }
    const raw = fs.readFileSync(srcPath, "utf8");
    const objData = JSON.parse(raw);
    const versionSetups = [];
    for (const prop of Object.getOwnPropertyNames(objData)) {
        versionSetups.push(new SetupEntry(prop, objData[prop]));
    }
    return versionSetups;
}