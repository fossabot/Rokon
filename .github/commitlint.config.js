// This file is used by commitlint to validate commit messages.
// It was adapted for monorepos
// While we're technically not a monorepo, it keeps things consistent
const fs = require('node:fs'),
  path = require('node:path'),
  { dirname: directoryName, basename: baseName } = require('node:path'),
  { lstatSync: fileInfo } = require('node:fs');

const readdirSync = (p, a = [], ignored = ['node_modules', '.git', '.next', '.husky']) => {
    if (fs.statSync(p).isDirectory()) {
      fs.readdirSync(p)
        .filter(f => {
          return baseName(f) !== ignored && !ignored?.includes(f);
        })

        .map(f => readdirSync(a[a.push(path.join(p, f)) - 1], a, ignored));
    }
    return a.filter(f => baseName(f) !== ignored && !ignored?.includes(f));
  },
  DEFAULT_SCOPES = [
    'commitlint',
    'sec',
    'security',
    'deps',
    'dependencies',
    'release',
    'actions',
    'docker',
    'windows',
    'linux',
    'macos'
  ],
  blacklistedScopes = new Set(['src', 'next', 'dist', 'out', 'main']),
  directoryNames = readdirSync('./')
    .map(file => directoryName(file))
    .map(entry => {
      const newEntry = fileInfo(entry);
      newEntry.name = baseName(entry);
      return newEntry;
    })
    .map(directory => directory.name)
    .map(s => {
      if (s.charAt(0) === '.') return s.slice(1);

      return s;
    })
    .map(s => {
      if (s.includes('-')) return s.split('-');

      return s;
    })
    .flat(Number.POSITIVE_INFINITY),
  scopes = [...new Set([...DEFAULT_SCOPES, ...directoryNames])]
    .map(s => {
      return s.replaceAll('_', '');
    })
    .filter(s => {
      return s.length > 0;
    })
    .filter(s => !blacklistedScopes.has(s));

console.log(scopes);
module.exports = {
  extends: ['@commitlint/config-conventional', 'monorepo'],
  rules: {
    'scope-enum': [2, 'always', scopes],
  },
};
