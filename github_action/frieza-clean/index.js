const core = require('@actions/core');

const setup = require('./lib/frieza');

(async () => {
  try {
    const access_key = core.getInput('access_key');
    const secret_key = core.getInput('secret_key');
    const region = core.getInput('region')
    const release = core.getInput('release');

    // Binary
    const pathToCLI = await setup.downloadBinary(release)
    core.debug(`Add ${pathToCLI} to PATH`)
    core.addPath(pathToCLI);

    // Credentials
    await setup.addCredentials(access_key, secret_key, region)

    // Snapshot
    await setup.makeSnapshot()

} catch (error) {
    core.setFailed(error.message);
}
})();