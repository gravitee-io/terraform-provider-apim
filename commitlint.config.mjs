export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'header-max-length': [2, 'always', 120],
    'type-enum': [
      2, 'always',
      ['feat', 'fix', 'chore', 'ci', 'docs', 'test', 'refactor', 'perf', 'build', 'revert', 'release'],
    ],
  },
};
