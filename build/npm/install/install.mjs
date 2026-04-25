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

export function findPackageBinDir(
  {
    packageName,
    packagesDir,
  },
) {
  if (!packagesDir || !packageName) {
    return null;
  }

  const packageRoot = path.isAbsolute(packagesDir) ? packagesDir : path.resolve(process.cwd(), packagesDir);
  const candidate = path.join(path.resolve(packageRoot), ...packageName.split('/'), 'bin');
  try {
    accessSync(candidate);
    return candidate;
  } catch {
    return null;
  }
}

function inferPackagesDir(scriptDir) {
  const candidates = [
    path.resolve(scriptDir, '..', 'packages'),
    path.resolve(scriptDir, '..', '..'),
  ];

  for (const candidate of candidates) {
    const markerPath = path.join(candidate, '@atlacp');
    try {
      accessSync(markerPath);
      return candidate;
    } catch {
      // continue
    }
  }

  return candidates[0];
}

export function resolveSourceBinDir(
  {
    packageName,
    packagesDir,
    scriptDir,
  },
) {
  if (!packageName) {
    throw new Error('Package name is required');
  }

  const configuredBinDir = findPackageBinDir({ packageName, packagesDir });
  if (configuredBinDir) {
    return configuredBinDir;
  }

  if (packagesDir) {
    throw new Error(`Could not find bin directory for package ${packageName} in ${packagesDir}`);
  }

  if (!scriptDir) {
    throw new Error(`Could not find bin directory for package ${packageName} (looked for: auto-detected script path)`);
  }

  const inferredPackagesDir = inferPackagesDir(scriptDir);
  const inferredBinDir = findPackageBinDir({ packageName, packagesDir: inferredPackagesDir });
  if (inferredBinDir) {
    return inferredBinDir;
  }

  throw new Error(`Could not find bin directory for package ${packageName} (looked for: ${inferredPackagesDir})`);
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
  console.log(`Installing atlacp tools for ${process.platform}/${process.arch}`);
  const packageName = detectPlatformPackage();
  if (!packageName) {
    throw new Error(`Unsupported platform/architecture: ${process.platform}/${process.arch}`);
  }
  console.log(`Detected platform package: ${packageName}`)

  const scriptDir = path.dirname(fileURLToPath(import.meta.url));
  const sourceBinDir = resolveSourceBinDir({
    packageName,
    packagesDir: process.env.ATLACP_PACKAGES_DIR,
    scriptDir,
  });
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
