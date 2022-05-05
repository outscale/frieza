const core = require('@actions/core');

const setup = require('./lib/frieza');

(async () => {
  try {
    await setup.cleanAccount()
} catch (error) {
    core.setFailed(error.message);
}
})();