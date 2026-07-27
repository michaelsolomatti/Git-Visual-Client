// git_visual.js
#!/usr/bin/env node
const simpleGit = require('simple-git');
const { program } = require('commander');
const chalk = require('chalk');

program
    .option('--graph', 'Show graph', true)
    .option('--all', 'Show all branches', false)
    .option('--decorate', 'Show refs', true)
    .option('--oneline', 'Compact output', false)
    .option('--author <name>', 'Filter by author')
    .option('--since <date>', 'Commits after date')
    .option('--until <date>', 'Commits before date')
    .option('-n, --max-count <number>', 'Limit commits', parseInt)
    .option('--color', 'Color output', true)
    .option('--path <path>', 'Repository path', '.')
    .parse(process.argv);

const opts = program.opts();
const git = simpleGit(opts.path);

async function getCommits() {
    const logOpts = {};
    if (opts.all) logOpts.all = true;
    if (opts.author) logOpts.author = opts.author;
    if (opts.since) logOpts.since = opts.since;
    if (opts.until) logOpts.until = opts.until;
    if (opts.maxCount) logOpts['-n'] = opts.maxCount;
    logOpts['--pretty'] = 'format:%h|%an|%ae|%ad|%s';
    const log = await git.log(logOpts);
    return log.all;
}

async function getRefs() {
    // Получаем ветки и теги
    const branches = await git.branch();
    const tags = await git.tags();
    const refs = {};
    for (const [name, branch] of Object.entries(branches.branches)) {
        const hash = branch.commit.substring(0, 7);
        if (!refs[hash]) refs[hash] = [];
        refs[hash].push(name);
    }
    for (const tag of tags.all) {
        // Упрощённо: получаем хеш для тега через show
        try {
            const show = await git.show(['-s', '--format=%H', tag]);
            const hash = show.substring(0, 7);
            if (!refs[hash]) refs[hash] = [];
            refs[hash].push(tag);
        } catch (e) {}
    }
    // HEAD
    const currentBranch = await git.branchLocal();
    const head = currentBranch.current;
    if (head) {
        const hash = currentBranch.branches[head].commit.substring(0, 7);
        if (!refs[hash]) refs[hash] = [];
        refs[hash].push('HEAD -> ' + head);
    }
    return refs;
}

function colorize(text, color) {
    if (!opts.color) return text;
    return chalk[color](text);
}

function printCommit(commit, prefix, last, refs) {
    const connector = last ? '└── ' : '├── ';
    const hashShort = commit.hash.substring(0, 7);
    const refStr = refs[hashShort] ? ' (' + refs[hashShort].join(', ') + ')' : '';
    if (opts.oneline) {
        console.log(`${prefix}${connector}${colorize(hashShort, 'yellow')}${colorize(refStr, 'blue')} ${colorize(commit.message, 'reset')}`);
    } else {
        console.log(`${prefix}${connector}${colorize(hashShort, 'yellow')}${colorize(refStr, 'blue')}`);
        console.log(`${prefix}    ${colorize('Author:', 'cyan')} ${commit.author_name} <${commit.author_email}>`);
        console.log(`${prefix}    ${colorize('Date:', 'green')} ${commit.date}`);
        console.log(`${prefix}    ${colorize(commit.message, 'reset')}\n`);
    }
}

async function main() {
    try {
        const commits = await getCommits();
        if (!commits.length) {
            console.log('No commits found.');
            return;
        }
        const refs = await getRefs();
        for (let i = 0; i < commits.length; i++) {
            const last = i === commits.length - 1;
            printCommit(commits[i], '', last, refs);
        }
    } catch (err) {
        console.error('Error:', err.message);
        process.exit(1);
    }
}

main();
