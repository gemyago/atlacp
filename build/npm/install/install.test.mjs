import assert from 'node:assert/strict';
import { mkdtemp, mkdir, readFile, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';

import {
  appendToPath,
  detectPlatformPackage,
  findPackageBinDir,
  resolveSourceBinDir,
  detectShellConfigFiles,
  isPathAlreadyConfigured,
} from './install.mjs';

test('detectPlatformPackage maps supported and unsupported platforms', () => {
  assert.equal(detectPlatformPackage('linux', 'x64'), '@atlacp/install-linux-amd64');
  assert.equal(detectPlatformPackage('linux', 'arm64'), '@atlacp/install-linux-arm64');
  assert.equal(detectPlatformPackage('darwin', 'x64'), '@atlacp/install-darwin-amd64');
  assert.equal(detectPlatformPackage('darwin', 'arm64'), '@atlacp/install-darwin-arm64');
  assert.equal(detectPlatformPackage('win32', 'x64'), null);
});

test('isPathAlreadyConfigured returns true when /.atlacp/bin is present', () => {
  const content = 'export PATH="$HOME/.atlacp/bin:$PATH"\n';
  assert.equal(isPathAlreadyConfigured(content, '/home/user/.atlacp/bin'), true);
});

test('isPathAlreadyConfigured returns false when /.atlacp/bin is absent', () => {
  const content = 'export PATH="$HOME/.local/bin:$PATH"\n';
  assert.equal(isPathAlreadyConfigured(content, '/home/user/.atlacp/bin'), false);
});

test('appendToPath appends export PATH line', async () => {
  const tempDir = await mkdtemp(path.join(tmpdir(), 'atlacp-install-test-'));
  const configFile = path.join(tempDir, '.zshrc');

  await writeFile(configFile, '# existing\n', 'utf8');
  await appendToPath(configFile, '/home/user/.atlacp/bin');

  const content = await readFile(configFile, 'utf8');
  assert.match(content, /export PATH="\$HOME\/.atlacp\/bin:\$PATH"/);
});

test('appendToPath is idempotent', async () => {
  const tempDir = await mkdtemp(path.join(tmpdir(), 'atlacp-install-test-'));
  const configFile = path.join(tempDir, '.bashrc');

  await writeFile(configFile, '# existing\n', 'utf8');
  await appendToPath(configFile, '/home/user/.atlacp/bin');
  await appendToPath(configFile, '/home/user/.atlacp/bin');

  const content = await readFile(configFile, 'utf8');
  const occurrences = content.split('export PATH="$HOME/.atlacp/bin:$PATH"').length - 1;

  assert.equal(occurrences, 1);
});

test('detectShellConfigFiles returns zsh config when shell ends with zsh', () => {
  const originalShell = process.env.SHELL;

  try {
    process.env.SHELL = '/bin/zsh';
    assert.deepEqual(detectShellConfigFiles(), ['~/.zshrc']);
  } finally {
    if (originalShell === undefined) {
      delete process.env.SHELL;
    } else {
      process.env.SHELL = originalShell;
    }
  }
});

test('detectShellConfigFiles returns bash config when shell ends with bash', () => {
  const originalShell = process.env.SHELL;

  try {
    process.env.SHELL = '/usr/local/bin/bash';
    assert.deepEqual(detectShellConfigFiles(), ['~/.bashrc']);
  } finally {
    if (originalShell === undefined) {
      delete process.env.SHELL;
    } else {
      process.env.SHELL = originalShell;
    }
  }
});

test('detectShellConfigFiles returns profile config for unknown shell', () => {
  const originalShell = process.env.SHELL;

  try {
    process.env.SHELL = '/usr/bin/fish';
    assert.deepEqual(detectShellConfigFiles(), ['~/.profile']);
  } finally {
    if (originalShell === undefined) {
      delete process.env.SHELL;
    } else {
      process.env.SHELL = originalShell;
    }
  }
});

test('resolveSourceBinDir resolves package bin from provided packages dir', async () => {
  const workspace = await mkdtemp(path.join(tmpdir(), 'atlacp-install-dist-'));
  const packagesBase = path.join(workspace, 'build', 'npm', 'packages');
  const localBinDir = path.join(packagesBase, '@atlacp', 'install-linux-amd64', 'bin');
  await mkdir(localBinDir, { recursive: true });
  await writeFile(path.join(localBinDir, 'bbcp'), '#!/bin/sh\necho local', 'utf8');

  const resolved = resolveSourceBinDir({ packageName: '@atlacp/install-linux-amd64', packagesDir: packagesBase });
  assert.equal(resolved, localBinDir);
});

test('findPackageBinDir returns null when packages dir missing', () => {
  const workspace = '/tmp/non-existing-atlacp-packages-root';
  const resolved = findPackageBinDir({
    packageName: '@atlacp/install-linux-amd64',
    packagesDir: path.join(workspace, 'build/npm/packages'),
  });
  assert.equal(resolved, null);
});

test('resolveSourceBinDir throws when package dir is missing', () => {
  const workspace = '/tmp/non-existing-atlacp-packages-root';

  assert.throws(
    () => resolveSourceBinDir({ packageName: '@atlacp/install-linux-amd64', packagesDir: workspace }),
    {
      message: /Could not find bin directory for package @atlacp/,
    },
  );
});

test('resolveSourceBinDir infers package dir from script location', async () => {
  const workspace = await mkdtemp(path.join(tmpdir(), 'atlacp-install-script-'));
  const scriptDir = path.join(workspace, 'build', 'npm', 'install');
  const packagesBase = path.join(workspace, 'build', 'npm', 'packages');
  const localBinDir = path.join(packagesBase, '@atlacp', 'install-linux-amd64', 'bin');
  await mkdir(localBinDir, { recursive: true });
  await writeFile(path.join(localBinDir, 'bbcp'), '#!/bin/sh\necho local', 'utf8');

  const resolved = resolveSourceBinDir({
    packageName: '@atlacp/install-linux-amd64',
    scriptDir,
  });
  assert.equal(resolved, localBinDir);
});
