import { accessSync } from 'node:fs';
import {
  chmod,
  copyFile,
  mkdir,
  readFile,
  readdir,
  writeFile,
} from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

const EXPORT_PATH_LINE = 'export PATH="$HOME/.atlacp/bin:$PATH"';

export function detectPlatformPackage(platform = process.platform, arch = process.arch) {
  const mapping = {
    'linux:x64': '@atlacp/install-linux-amd64',
    'linux:arm64': '@atlacp/install-linux-arm64',
    'darwin:x64': '@atlacp/install-darwin-amd64',
    'darwin:arm64': '@atlacp/install-darwin-arm64',
  };

  return mapping[`${platform}:${arch}`] ?? null;
}

export function findPlatformBinDir(packageName, baseDir = path.dirname(fileURLToPath(import.meta.url))) {
  if (!packageName) {
    throw new Error('Platform package name is required');
  }

  let currentDir = path.resolve(baseDir);
  const packagePathParts = packageName.split('/');

  while (true) {
    const binDir = path.join(currentDir, 'node_modules', ...packagePathParts, 'bin');

    try {
      accessSync(binDir);
      return binDir;
    } catch {
      // continue searching parent dirs
    }

    const parentDir = path.dirname(currentDir);
    if (parentDir === currentDir) {
      break;
    }

    currentDir = parentDir;
  }

  throw new Error(`Could not find bin directory for package ${packageName}`);
}

export async function ensureDir(dirPath) {
  await mkdir(dirPath, { recursive: true });
}

export async function copyBinaries(srcDir, destDir) {
  await ensureDir(destDir);

  const entries = await readdir(srcDir, { withFileTypes: true });
  const files = entries.filter((entry) => entry.isFile());

  for (const file of files) {
    const srcFile = path.join(srcDir, file.name);
    const destFile = path.join(destDir, file.name);

    await copyFile(srcFile, destFile);
    await chmod(destFile, 0o755);
  }
}

export function detectShellConfigFiles() {
  const shell = process.env.SHELL ?? '';

  if (shell.endsWith('zsh')) {
    return ['~/.zshrc'];
  }

  if (shell.endsWith('bash')) {
    return ['~/.bashrc'];
  }

  return ['~/.profile'];
}

export function isPathAlreadyConfigured(configFileContent, binDir) {
  if (!configFileContent) {
    return false;
  }

  return configFileContent.includes('/.atlacp/bin') || configFileContent.includes(binDir);
}

function expandHome(filePath) {
  if (!filePath.startsWith('~/')) {
    return filePath;
  }

  return path.join(os.homedir(), filePath.slice(2));
}

export async function appendToPath(configFile, binDir) {
  const configFilePath = expandHome(configFile);

  let existingContent = '';
  try {
    existingContent = await readFile(configFilePath, 'utf8');
  } catch (err) {
    if (err?.code !== 'ENOENT') {
      throw err;
    }
  }

  if (isPathAlreadyConfigured(existingContent, binDir)) {
    return false;
  }

  const separator = existingContent.length > 0 && !existingContent.endsWith('\n') ? '\n' : '';
  const newContent = `${existingContent}${separator}${EXPORT_PATH_LINE}\n`;

  await ensureDir(path.dirname(configFilePath));
  await writeFile(configFilePath, newContent, 'utf8');

  return true;
}

export async function run() {
  const packageName = detectPlatformPackage();
  if (!packageName) {
    throw new Error(`Unsupported platform/architecture: ${process.platform}/${process.arch}`);
  }

  const scriptDir = path.dirname(fileURLToPath(import.meta.url));
  const sourceBinDir = findPlatformBinDir(packageName, scriptDir);
  const installBaseDir = path.join(os.homedir(), '.atlacp');
  const destinationBinDir = path.join(installBaseDir, 'bin');

  await ensureDir(destinationBinDir);
  await copyBinaries(sourceBinDir, destinationBinDir);

  const shellConfigFiles = detectShellConfigFiles();
  for (const shellConfigFile of shellConfigFiles) {
    await appendToPath(shellConfigFile, destinationBinDir);
  }

  console.log('atlacp binaries installed to ~/.atlacp/bin');
  console.log('Restart your shell or run: source ~/.profile (or your shell config file)');
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  run().catch((err) => {
    console.error(`Failed to install atlacp binaries: ${err.message}`);
    process.exitCode = 1;
  });
}
