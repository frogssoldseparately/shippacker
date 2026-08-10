import fs from "fs";
import path from "path";
import { XmlResource } from "./XmlResource.js";
import { SetupEntry } from "./SetupEntry.js";

/**
 * Gets version and platform information from a filename.
 * @param {string} filePath - Filename being parsed.
 * @returns [Version string, Platform string]
 */
function getTargetInfo(filePath) {
    const parts = path.basename(filePath).split("_");
    const version = parts.slice(0, 2).join("_");
    const platform = parts.slice(2, parts.length - 1).join("_");
    return [version, platform];
}

/**
 * Find and return any sample name mappings in a directory.
 * @param {string} dir - Directory being checked for mappings
 * @returns Map whose keys are version-platform pairs and values are replacement
 * mappings of the corresponding ROM target.
 */
function findReplacements(dir) {
    const sourceFiles = fs
        .readdirSync(dir)
        .filter(elem => elem.endsWith("Audio.txt"))
        .map(elem => path.join(dir, elem));
    const out = new Map();
    for (const filePath of sourceFiles) {
        const [version, platform] = getTargetInfo(filePath);
        const key = `${version}_${platform}`;
        if (out.has(key)) {
            console.log(`Replacements:Duplicate: disregarding "${filePath}"`);
        } else {
            console.log(`Replacements:Found: ${filePath}`);
            const raw = fs.readFileSync(filePath, "utf8");
            const replacements = raw
                .split('\n')
                .map(elem => elem.trim());
            out.set(key, replacements);
        }
    }
    return out;
}

/**
 * Generates abridged Audio.xml files whose names match what is found in mm.o2r.
 * @param {string} source - Directory with unmodified xmls.
 * @param {string} destination - Directory where modified xmls will be placed.
 * @returns Array of names of files written.
 */
export function makeStubs(source, destination) {
    const writtenFiles = [];
    const sourceFiles = fs
        .readdirSync(source)
        .filter(elem => elem.endsWith("Audio.xml"))
        .map(elem => path.join(source, elem));
    const replacementMap = findReplacements(source);
    for (const resourcePath of sourceFiles) {
        const [version, platform] = getTargetInfo(resourcePath);
        // Try to get a replacement mappings that works for all platforms of
        // this version. Otherwise, get a platform specific replacement.
        const replacements = replacementMap.get(`${version}_All`)
            ?? replacementMap.get(`${version}_${platform}`);
        if (replacements == null) {
            console.log(`Replacements:Missing: For "${
                resourcePath
            }"`);
            continue;
        }
        const outputName = `${version}_${platform}_Audio_Stub.xml`;
        const outputPath = path.join(destination, outputName);
        try {
            // Wrapped in a class in the case scenario in needs to do more
            // in the future.
            (new XmlResource(resourcePath, replacements, outputPath))
                .writeStub();
            console.log(`Stub:Good: Wrote ${outputPath}`);
            writtenFiles.push(outputName);
        } catch (e) {
            console.log(`Stub:Bad: Could not write to "${
                outputPath
            }". Reason:\n\t${
                e.message
            }`);
        }
    }
    return writtenFiles;
}

/**
 * Generates audio.go with resource embeds and a map for access in Ship Packer.
 * @param {string} srcPath - audio.template path.
 * @param {string} destPath - audio.go path.
 * @param {Array<string>} writtenFiles - Names of Audio.xml stubs written.
 * @param {Array<SetupEntry>} versionSetups - Setup information for each
 * supported version.
 */
export function fillAudioTemplate(srcPath, destPath, writtenFiles, 
        versionSetups) {
    // Generate template fillers
    let embedEntries = "";
    let mapEntries = "";
    // let versionNames = "";
    // let versionNameMap = new Map();
    for (const filename of writtenFiles) {
        const filenameParts = filename.split('_')
        const keyname = filenameParts.slice(0, 4).join('_').toLowerCase();
        const varname = keyname + '_audio_file';
        embedEntries += `\n//go:embed ${filename}\nvar ${varname} []byte`;
        mapEntries += `\n\t"${keyname}" : ${varname},`;
    }
    let versionNames = "";
    for (const setup of versionSetups) {
        versionNames += `\n\t"${setup.name}",`;
    } 
    const injections = new Map([
        ["xml_embeds", embedEntries.slice(1)], // remove leading \n
        ["map_entries", mapEntries.slice(2)], // remove leading \n\t
        ["version_names", versionNames.slice(2)], // remove leading \n\t
    ]);
    fillTemplate(srcPath, destPath, injections);
}

/**
 * Generates setup.go with global config setups for each supported version.
 * @param {string} srcPath - setup.template path.
 * @param {string} destPath - setup.go path.
 * @param {Array<SetupEntry>} versionSetups - Setup information for each
 * supported version.
 */
export function fillSetupTemplate(srcPath, destPath, versionSetups) {
    let setupCases = "";
    for (const setup of versionSetups) {
        setupCases += `\n\tcase "${setup.name}":\n\t\t`;
        let bodyString = "";
        for (const assignment of setup.assignStrings) {
            bodyString += `${assignment}\n\t\t`
        }
        setupCases += bodyString.slice(0, -3); // remove trailing \n\t\t
    }
    const injections = new Map([
        ["setup_cases", setupCases.slice(2)] // remove leading \n\t
    ]);
    fillTemplate(srcPath, destPath, injections);
}

/**
 * Populates a template and writes it to the destination path.
 * @param {string} srcPath - template path.
 * @param {string} destPath - destination path.
 * @param {Map<string, string>} injections - replacements in the template.
 */
function fillTemplate(srcPath, destPath, injections) {
    let content = fs.readFileSync(srcPath, "utf8");
    for (const [key, value] of injections) {
        content = content.replaceAll(`##:${key}:##`, value);
    }
    fs.writeFileSync(destPath, content, "utf8");
    console.log(`Template:Good: Wrote ${destPath}`);
}