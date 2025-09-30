#!/usr/bin/env node
'use strict';

/**
 * Inserts the following snippet just before the first occurrence of:
 *   providerHTTPTransportOpts := ProviderHTTPTransportOpts{
 *
 * Snippet:
 *   // BEGIN GRAVITEE CLOUD INIT
 *   serverUrl = CloudInitializer(&security, serverUrl, &data, resp)
 *   // END GRAVITEE CLOUD INIT
 *
 * Features:
 * - Zero dependencies
 * - Idempotent (checks for BEGIN marker)
 * - Preserves file's newline style and indentation level of the anchor line
 */

const fs = require('node:fs');
const fsp = require('node:fs/promises');
const path = require('node:path');

const DEFAULT_TARGET = 'internal/provider/provider.go';
const ANCHOR = 'providerHTTPTransportOpts := ProviderHTTPTransportOpts{';
const SNIPPET_LINES = [
    '// BEGIN GRAVITEE CLOUD INIT',
    'serverUrl = CloudInitializer(security, serverUrl, &data, resp)',
    '// END GRAVITEE CLOUD INIT',
];

function parseArgs(argv) {
    const args = {};
    for (let i = 2; i < argv.length; i++) {
        const a = argv[i];
        if (a === '--path' && argv[i + 1]) {
            args.path = argv[++i];
        }
    }
    return args;
}

function detectNewline(text) {
    return text.includes('\r\n') ? '\r\n' : '\n';
}

function indentSnippet(snippetLines, indent, nl) {
    // Keep blank lines blank; otherwise prefix with indent
    return snippetLines
        .map((line) => (line.length ? indent + line : line))
        .join(nl);
}

async function main() {
    const { path: targetPathArg } = parseArgs(process.argv);
    const targetPath = targetPathArg ? targetPathArg : DEFAULT_TARGET;

    const absPath = path.resolve(process.cwd(), targetPath);

    if (!fs.existsSync(absPath)) {
        console.error(`Error: File not found: ${absPath}`);
        process.exitCode = 1;
        return;
    }

    let content;
    try {
        content = await fsp.readFile(absPath, 'utf8');
    } catch (e) {
        console.error(`Error: Failed to read file: ${absPath}\n${e.message}`);
        process.exitCode = 1;
        return;
    }

    // Idempotency check
    if (content.includes('// BEGIN GRAVITEE CLOUD INIT')) {
        console.log('Info: Snippet already present. No changes made.');
        return;
    }

    const idx = content.indexOf(ANCHOR);
    if (idx === -1) {
        console.error(`Error: Anchor not found in file: "${ANCHOR}"`);
        process.exitCode = 1;
        return;
    }

    const nl = detectNewline(content);

    // Find start of the anchor line to insert before it
    const lineStart = content.lastIndexOf('\n', idx - 1);
    const anchorLineStart = lineStart === -1 ? 0 : lineStart + 1;
    const lineEnd = content.indexOf('\n', idx);
    const anchorLine = content.slice(anchorLineStart, lineEnd === -1 ? content.length : lineEnd);

    const indentMatch = anchorLine.match(/^\s*/);
    const indent = indentMatch ? indentMatch[0] : '';

    const snippet = indentSnippet(SNIPPET_LINES, indent, nl) + nl;

    const before = content.slice(0, anchorLineStart);
    const after = content.slice(anchorLineStart);
    const updated = before + snippet + after;

    try {
        await fsp.writeFile(absPath, updated, 'utf8');
        console.log(`Success: Cloud initializer code snippet inserted in ${path.relative(process.cwd(), absPath)}`);
    } catch (e) {
        console.error(`Error: Failed to write updated file: ${e.message}`);
        process.exitCode = 1;
    }
}

main();