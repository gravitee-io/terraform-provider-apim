/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const { log } = console;

const newLoggerFn =
  (color) =>
  (...args) =>
    log(color(...args));

const $Quote = $.quote;
const $NoQuote = (unescaped) => unescaped;

// Color loggers
const green = newLoggerFn(chalk.green);
const blue = newLoggerFn(chalk.blue);
const magenta = newLoggerFn(chalk.magenta);
const yellow = newLoggerFn(chalk.yellow);
const red = newLoggerFn(chalk.red);

export const LOG = Object.seal({ green, blue, magenta, yellow, red, log });

// Record time taken by a fn to execute
export async function time(fn) {
  const start = Date.now();
  await fn();
  const end = Date.now();
  green(`Done in ${(end - start) / 1000}s`);
}

// Toggle zx verbosity (print commands before executing them)
export function toggleVerbosity(verbose = false) {
  $.verbose = verbose;
}

// Disable quote escaping for zx
export function setNoQuoteEscape() {
  $.quote = $NoQuote;
}

// Enable quote escaping for zx
export function setQuoteEscape() {
  $.quote = $Quote;
}

// Path to the project root directory
export const PROJECT_DIR = path.join(__dirname, "..", "..");

export function isNonEmptyString(str) {
  return String(str) === str && str.trim().length > 0;
}

export function isEmptyString(str) {
  return !isNonEmptyString(str);
}