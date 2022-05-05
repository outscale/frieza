const core = require('@actions/core');
const tc = require('@actions/tool-cache');
const io = require('@actions/io');
const exec = require('@actions/exec');

const fetch = require('node-fetch')
const os = require('os');
const path = require('path');

const default_profile_name = 'action'
const default_snapshot_name = 'snapshot-action'

async function getRelease(release) {
    let url = ''
    if (release == '') {
        url = "https://api.github.com/repos/outscale-dev/frieza/releases/latest"
    } else {
        url = `https://api.github.com/repos/outscale-dev/frieza/releases/${release}`
    }
    let response = await fetch(url);
    let data = await response.json();
    return data;
}

function getAssetURL(data, asset_name) {
    for (const element of data) {
        if (element["name"] == asset_name) {
            return element["browser_download_url"]
        }
    }
    return ""
}

async function copyBinary(pathToCLI, release_tag) {
    const exeSuffix = os.platform().startsWith('win') ? '.exe' : '';

    try {
        source = [pathToCLI, `frieza_${release_tag}${exeSuffix}`].join(path.sep);
        target = [pathToCLI, `frieza${exeSuffix}`].join(path.sep);
        core.debug(`Moving ${source} to ${target}.`);
        await io.mv(source, target);
    } catch (e) {
        core.error(`Unable to move ${source} to ${target}.`);
        throw e;
    }
}

async function downloadBinary(release) {
    let release_data = await getRelease(release)

    let release_tag = release_data["tag_name"]
    if (release_tag.startsWith("v")) {
        release_tag = release_tag.substring(1)
    }

    const asset_name = "frieza_" + release_tag + "_" + mapOS(os.platform()) + "_" + mapArch(os.arch()) + ".zip"
    const url = getAssetURL(release_data["assets"], asset_name)

    core.debug(`Downloading Frieza from ${url}`);
    const downloadedPath = await tc.downloadTool(url)

    core.debug('Extracting Frieza zip file');
    const pathToCLI = await tc.extractZip(downloadedPath);
    core.debug(`Frieza path is ${pathToCLI}.`);

    if (!downloadedPath || !pathToCLI) {
        throw new Error(`Unable to download Frieza from ${url}`);
    }

    await copyBinary(pathToCLI, `v${release_tag}`)

    return pathToCLI;
}

function mapArch(arch) {
    const mappings = {
        'x32': '386',
        'x64': 'amd64'
    }
    return mappings[arch] || arch
}

function mapOS(os) {
    const mappings = {
        'win32': 'windows'
    };
    return mappings[os] || os;
}

async function addCredentials(access_key, secret_key, region) {
    core.debug(`Add credentials to frieza`);
    await exec.exec('frieza', ['profile', 'new', 'outscale_oapi', `--region=${region}`, `--ak=${access_key}`, `--sk=${secret_key}`, default_profile_name]);
}

async function makeSnapshot() {
    core.debug(`Make a snapshot`);
    await exec.exec('frieza', ['snapshot', 'new', default_snapshot_name, default_profile_name]);
}

async function cleanAccount() {
    core.debug(`Clean account`);
    await exec.exec('frieza', ['clean', '--auto-approve', default_snapshot_name]);
}


exports.makeSnapshot = makeSnapshot;
exports.addCredentials = addCredentials;
exports.downloadBinary = downloadBinary;
exports.cleanAccount = cleanAccount;