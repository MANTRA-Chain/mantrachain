import type { Config } from 'jest'
import { defaults } from 'jest-config'

const config: Config = {
  reporters: ['default'],
  cacheDirectory: '.jest/cache',
  coverageDirectory: '.jest/coverage',
  bail: true,
  globalSetup: './globalSetup.ts',
  testTimeout: 600000,
  moduleFileExtensions: [...defaults.moduleFileExtensions, 'cjs'],
  setupFilesAfterEnv: ['jest-extended/all'],
  transform: {
    '^.+\\.[t|j]sx?$': 'babel-jest',
  },
}

export default config
