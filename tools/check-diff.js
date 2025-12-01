import { execFile as cbExecFile } from 'node:child_process';
import { promisify } from 'node:util';

const execFile = promisify(cbExecFile);

/**
 *
 * @param {string} msg
 * @param {string} file
 */
const emitError = (msg, file) => {
  process.stdout.write(`::error file=${file}::${msg}` + '\n');
};

/**
 *
 * @param {string} file
 */
const reportError = (file) => {
  emitError('Unexpected difference is occurred', file);
};

const main = async () => {
  const { stdout } = await execFile('git', ['diff', '--name-only']);
  const modifiedFiles = stdout
    .trim()
    .split('\n')
    .filter((file) => file !== '');
  for (const file of modifiedFiles) {
    reportError(file);
  }
  if (modifiedFiles.length > 0) {
    process.exit(1);
  }
};

await main();
