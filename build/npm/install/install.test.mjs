import assert from 'node:assert/strict';
import { mkdtemp, readFile, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';
import test from 'node:test';

import {
  appendToPath,
  detectPlatformPackage,
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
